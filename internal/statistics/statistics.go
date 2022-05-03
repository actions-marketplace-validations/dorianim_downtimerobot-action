package statistics

import (
	"fmt"
	"math"
	"time"

	"github.com/dorianim/downtimerobot/internal/crawler"
)

// When chaging, also update types in statistics.ts

type Service struct {
	Name            string           `json:"name"`
	Host            string           `json:"host"`
	Type            string           `json:"type"`
	Up              bool             `json:"up"`
	Disabled        bool             `json:"disabled"`
	Uptime          UptimeStatistics `json:"uptime"`
	DailyStatistics [90]float32      `json:"dailyStatistics"`
	logs            []ServiceLog
	responseTimes   []ServiceResponseTime
}

type DetailedService struct {
	Service
	Logs          []ServiceLog          `json:"logs"`
	ResponseTimes []ServiceResponseTime `json:"responseTimes"`
}

type ServiceLog struct {
	Up             bool   `json:"up"`
	Disabled       bool   `json:"disabled"`
	TimeString     string `json:"timeString"`
	DurationString string `json:"durationString"`
	Status         struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
}

type ServiceResponseTime struct {
	Value      int64  `json:"value"`
	TimeString string `json:"timeString"`
}

type ServiceList struct {
	Services   []Service  `json:"services"`
	Days       [90]string `json:"days"`
	Statistics Statistics `json:"statistics"`
	TimeZone   string     `json:"timeZone"`
}

type ServiceDetails struct {
	Service  DetailedService `json:"service"`
	Days     [90]string      `json:"days"`
	TimeZone string          `json:"timeZone"`
}

type UptimeStatistics struct {
	OneDay     float32 `json:"1"`
	SevenDays  float32 `json:"7"`
	ThirtyDays float32 `json:"30"`
	NinetyDays float32 `json:"90"`
}

type CountStatistics struct {
	Up       int `json:"up"`
	Down     int `json:"down"`
	Disabled int `json:"disabled"`
	Total    int `json:"total"`
}

type Statistics struct {
	Uptime UptimeStatistics `json:"uptime"`
	Counts CountStatistics  `json:"counts"`
}

func Generate(crawledServices []crawler.Service) (ServiceList, []ServiceDetails, error) {
	dayStrings := generateDayStrings()
	timeZone := time.Now().Format("-07:00")

	serviceList := generateServiceList(crawledServices, dayStrings, timeZone)
	serviceDetailsList := generateServiceDetails(serviceList.Services, dayStrings, timeZone)

	return serviceList, serviceDetailsList, nil
}

func generateServiceList(crawledServices []crawler.Service, dayStrings [90]string, timeZone string) ServiceList {
	serviceList := ServiceList{}
	serviceList.Days = dayStrings
	serviceList.TimeZone = timeZone
	serviceList.Services = calculateAllServiceStatistics(crawledServices)
	serviceList.Statistics = calculateStatistics(serviceList.Services)
	return serviceList
}

func generateServiceDetails(services []Service, dayStrings [90]string, timeZone string) []ServiceDetails {
	serviceDetailsList := make([]ServiceDetails, 0)
	for _, service := range services {
		serviceDetails := ServiceDetails{}
		serviceDetails.Days = dayStrings
		serviceDetails.TimeZone = timeZone
		serviceDetails.Service = DetailedService{service, service.logs, service.responseTimes}
		serviceDetailsList = append(serviceDetailsList, serviceDetails)
	}
	return serviceDetailsList
}

// == Service statistics ==
func calculateAllServiceStatistics(crawledServices []crawler.Service) []Service {
	services := make([]Service, len(crawledServices))
	for i, crawledService := range crawledServices {
		services[i] = Service{}
		services[i].Name = crawledService.GetName()
		services[i].Host = crawledService.GetHost()
		services[i].Disabled = crawledService.IsDisabled()
		services[i].Type = crawledService.GetType()
		services[i].DailyStatistics = calculateServiceStatistics(crawledService)
		services[i].Up = crawledService.IsUp()
		services[i].Uptime = calculateServiceUptimeStatistics(services[i])
		services[i].responseTimes = getServiceResponseTimes(crawledService)
		services[i].logs = generateServiceLogs(crawledService)
	}
	return services
}

func calculateServiceStatistics(crawledService crawler.Service) [90]float32 {
	var result [90]float32
	now := time.Now()
	tmpDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	for i := 0; i < 90; i++ {
		startDate := tmpDate.AddDate(0, 0, -i)
		endDate := startDate.Add(time.Hour * 24)
		result[i] = calculateUptime(getDataPointsBetween(
			startDate.Unix(),
			endDate.Unix(),
			crawledService.GetHistoricData(),
		))
	}
	return result
}

func calculateUptime(dataPoints []crawler.HistoricDataPoint) float32 {
	if len(dataPoints) <= 0 {
		return -1
	}

	var sum float32 = 0.0
	var count int = 0
	for _, datadataPoint := range dataPoints {
		if !datadataPoint.IsDisabled() {
			count++
		}
		if datadataPoint.IsUp() {
			sum++
		}
	}

	if count == 0 {
		return -1
	}
	return round(sum / float32(count))
}

func calculateServiceUptimeStatistics(service Service) UptimeStatistics {
	uptime := UptimeStatistics{}
	var sum float32 = 0.0
	var count float32 = 0.0
	for i := 0; i < 90; i++ {
		if service.DailyStatistics[i] >= 0.0 {
			sum += service.DailyStatistics[i]
			count++
		}

		var average float32 = -1
		if count > 0 {
			average = round(sum / count)
		}
		if i == 0 {
			uptime.OneDay = average
		} else if i == 6 {
			uptime.SevenDays = average
		} else if i == 29 {
			uptime.ThirtyDays = average
		} else if i == 89 {
			uptime.NinetyDays = average
			break
		}
	}
	return uptime
}

// TODO: make configurable?
// getServiceResponseTimes returns the services response times of the last 24 hours
func getServiceResponseTimes(crawledService crawler.Service) []ServiceResponseTime {
	responseTimes := make([]ServiceResponseTime, 0)
	to := time.Now()
	from := to.Add(-time.Hour * 24)
	dataPoints := getDataPointsBetween(from.Unix(), to.Unix(), crawledService.GetHistoricData())
	for _, dataPoint := range dataPoints {
		value := dataPoint.GetResponseTime()
		if value <= 0 {
			continue
		}
		dateTime := time.Unix(dataPoint.GetTimestamp(), 0).Format("January 02, 2006, 14:04")
		responseTimes = append(responseTimes, ServiceResponseTime{value, dateTime})
	}
	return responseTimes
}

func generateServiceLogs(crawledService crawler.Service) []ServiceLog {
	logs := make([]ServiceLog, 0)
	to := time.Now()
	from := to.AddDate(0, 0, -90)
	dataPoints := getDataPointsBetween(from.Unix(), to.Unix(), crawledService.GetHistoricData())

	var previousStatusCode int = -1
	var previousLogTimestamp time.Time
	var previousServiceLog ServiceLog

	for _, dataPoint := range dataPoints {
		generateServiceLog(dataPoint, &logs, &previousStatusCode, &previousLogTimestamp, &previousServiceLog)
	}

	if len(dataPoints) > 0 {
		generateServiceLog(dataPoints[len(dataPoints)-1], &logs, nil, &previousLogTimestamp, &previousServiceLog)
	}

	return logs
}

// TODO: refactor, this is super ugly lol
// generateServiceLog generates a log for one datapoint if it is different from the one before. If previousStatusCode is nil, it will also generate the log
func generateServiceLog(dataPoint crawler.HistoricDataPoint, logs *[]ServiceLog, previousStatusCode *int, previousLogTimestamp *time.Time, previousServiceLog *ServiceLog) {
	if previousStatusCode != nil && dataPoint.GetStatusCode() == *previousStatusCode {
		return
	}

	logTime := time.Unix(dataPoint.GetTimestamp(), 0)

	serviceLog := ServiceLog{}
	serviceLog.Up = dataPoint.IsUp()
	serviceLog.Disabled = dataPoint.IsDisabled()
	serviceLog.TimeString = logTime.Format("January 02, 2006, 15:04")
	serviceLog.Status.Code = dataPoint.GetStatusCode()
	serviceLog.Status.Message = dataPoint.GetStatusMessage()

	if previousStatusCode == nil || *previousStatusCode != -1 {
		previousServiceLog.DurationString = durationAsString(*previousLogTimestamp, logTime)
		*logs = append(*logs, *previousServiceLog)
	}

	if previousStatusCode != nil {
		*previousStatusCode = dataPoint.GetStatusCode()
		*previousLogTimestamp = logTime
		*previousServiceLog = serviceLog
	}
}

// == General statistics ==
func calculateStatistics(services []Service) Statistics {
	statistics := Statistics{}
	statistics.Counts = calculateCountStatistics(services)
	statistics.Uptime = calculateUptimeStatistics(services, statistics.Counts)
	return statistics
}

func calculateCountStatistics(services []Service) CountStatistics {
	counts := CountStatistics{}
	for _, service := range services {
		if service.Disabled {
			counts.Disabled++
		} else if service.Up {
			counts.Up++
		} else {
			counts.Down++
		}
		counts.Total++
	}
	return counts
}

func calculateUptimeStatistics(services []Service, counts CountStatistics) UptimeStatistics {
	uptime := UptimeStatistics{}
	totalServices := float32(counts.Up + counts.Down)
	for _, service := range services {
		if service.Disabled {
			continue
		}

		uptime.OneDay += round(service.Uptime.OneDay / totalServices)
		uptime.SevenDays += round(service.Uptime.SevenDays / totalServices)
		uptime.ThirtyDays += round(service.Uptime.ThirtyDays / totalServices)
		uptime.NinetyDays += round(service.Uptime.NinetyDays / totalServices)
	}
	return uptime
}

// == Helpers ==

func durationAsString(from time.Time, to time.Time) string {
	durationHours := (to.Unix() - from.Unix()) / (60 * 60)
	durationMinutes := ((to.Unix() - from.Unix()) % (60 * 60)) / 60
	return fmt.Sprintf("%d h, %d min", durationHours, durationMinutes)
}

func round(num float32) float32 {
	return float32(math.Round(float64(num*100000)) / 100000)
}

// getDataPointsBetween assumes that the datapoints are aleardy sorted
func getDataPointsBetween(startTimestamp int64, endTimestamp int64, dataPoints []crawler.HistoricDataPoint) []crawler.HistoricDataPoint {
	result := make([]crawler.HistoricDataPoint, 0)
	for _, dataPoint := range dataPoints {
		if dataPoint.GetTimestamp() >= startTimestamp && dataPoint.GetTimestamp() <= endTimestamp {
			result = append(result, dataPoint)
		}
	}
	return result
}

func generateDayStrings() [90]string {
	var dayStrings [90]string
	tmpDate := time.Now()
	for i := 0; i < 90; i++ {
		dayStrings[i] = tmpDate.AddDate(0, 0, -i).Format("January 02, 2006")
	}
	return dayStrings
}
