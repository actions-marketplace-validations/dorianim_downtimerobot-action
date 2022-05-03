package crawler

import (
	"io/ioutil"
	"os"
	"reflect"

	"github.com/goccy/go-json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Service describes a service wich is monitored by downtimerobot
type Service interface {
	crawl() HistoricDataPoint
	setHistoricData([]rawHistoricDataPoint)

	GetHost() string
	GetName() string
	GetType() string
	IsDisabled() bool
	IsUp() bool
	GetHistoricData() []HistoricDataPoint
}

type genericService struct {
	Service
	Name         string `json:"name"`
	Host         string `json:"host"`
	Disabled     bool   `json:"disabled"`
	historicData []HistoricDataPoint
}

// HistoricDataPoint is the status of a service at a certain point of time
type HistoricDataPoint interface {
	IsUp() bool
	IsDisabled() bool
	GetStatusCode() int
	GetStatusMessage() string
	GetResponseTime() int64
	GetTimestamp() int64

	getRawDataPoint() rawHistoricDataPoint
}

type config struct {
	Services struct {
		HTTPS   []*httpsService
		Ping    []*PingService
		Port    []*PortService
		Pattern []*patternService
	}
}

type rawHistoricDataPoint struct {
	Timestamp    int64 `json:"t"`
	StatusCode   int   `json:"c"`
	ResponseTime int64 `json:"r"`

	// StatusMessage is only used for non-standard errors
	// it may be empty in many cases. Use Service::GetStatus
	StatusMessage string `json:"m"`
}

type rawHistoricData map[string][]rawHistoricDataPoint

const historicDataFile = "./historicData.json"

// CrawlServices crawls all configured services and adds the result to their historic data
func CrawlServices() ([]Service, error) {
	services, err := LoadServices()
	if err != nil {
		return nil, err
	}

	crawlServices(services)
	if err := storeHistoricData(services); err != nil {
		return nil, err
	}

	return services, err
}

func LoadServices() ([]Service, error) {
	conf, err := loadConfig()
	if err != nil {
		return nil, err
	}

	historicData, err := loadHistoricData()
	if err != nil {
		return nil, err
	}

	return loadServices(conf, historicData), nil
}

func loadConfig() (*config, error) {
	conf := &config{}
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func loadServices(conf *config, historicData rawHistoricData) []Service {
	result := make([]Service, 0)

	serviceTypes := reflect.ValueOf(conf.Services)
	for i := 0; i < serviceTypes.NumField(); i++ {
		field := serviceTypes.Field(i)
		for j := 0; j < field.Len(); j++ {
			service := field.Index(j).Interface().(Service)
			injectHistoricDataIntoService(historicData, service)
			result = append(result, service)
		}
	}

	return result
}

func crawlServices(services []Service) {
	for _, service := range services {
		crawlService(service)
	}
}

func dataPointsToRawDataPoints(dataPoints []HistoricDataPoint) []rawHistoricDataPoint {
	result := make([]rawHistoricDataPoint, len(dataPoints))
	for i, d := range dataPoints {
		result[i] = d.getRawDataPoint()
	}
	return result
}

func loadHistoricData() (rawHistoricData, error) {
	jsonFile, err := os.Open(historicDataFile)
	defer jsonFile.Close()

	if err != nil && os.IsNotExist(err) {
		return rawHistoricData{}, nil
	} else if err != nil {
		return nil, err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	data := make(rawHistoricData)
	err = json.Unmarshal(byteValue, &data)
	return data, err
}

func storeHistoricData(services []Service) error {
	rawData := rawHistoricData{}
	for _, service := range services {
		rawData[service.GetHost()] = dataPointsToRawDataPoints(service.GetHistoricData())
	}

	data, _ := json.MarshalIndent(rawData, "", " ")
	return ioutil.WriteFile(historicDataFile, data, 0644)
}

func injectHistoricDataIntoService(data rawHistoricData, service Service) {
	if serviceData, ok := data[service.GetHost()]; ok {
		service.setHistoricData(serviceData)
	}
}

func crawlService(service Service) {
	newDataPoint := service.crawl()
	if newDataPoint.IsDisabled() {
		log.WithFields(log.Fields{
			"service": service.GetHost(),
			"type":    service.GetType(),
		}).Info("Service is DISABLED")
	} else if newDataPoint.IsUp() {
		log.WithFields(log.Fields{
			"service": service.GetHost(),
			"type":    service.GetType(),
		}).Info("Service is UP")
	} else {
		log.WithFields(log.Fields{
			"service":       service.GetHost(),
			"type":          service.GetType(),
			"statusMessage": newDataPoint.GetStatusMessage(),
			"statusCode":    newDataPoint.GetStatusCode(),
		}).Warn("Service is Down")
	}
}

// == genericService ==
func (service *genericService) GetHost() string {
	return service.Host
}

func (service *genericService) GetName() string {
	return service.Name
}

func (service *genericService) IsDisabled() bool {
	return service.Disabled
}

func (genericService *genericService) IsUp() bool {
	historicData := genericService.GetHistoricData()
	historicDataLength := len(historicData)
	if historicDataLength > 0 {
		return historicData[historicDataLength-1].IsUp()
	}
	return false
}

func (genericService *genericService) GetHistoricData() []HistoricDataPoint {
	return genericService.historicData
}

func (genericService *genericService) GetType() string {
	return "generic"
}

func (genericService *genericService) crawl() HistoricDataPoint {
	return nil
}

func (genericService *genericService) setHistoricData([]rawHistoricDataPoint) {

}

// == rawHistoricDataPoint ==

func (dataPoint rawHistoricDataPoint) IsDisabled() bool {
	return dataPoint.StatusCode == -1
}

func (dataPoint rawHistoricDataPoint) GetResponseTime() int64 {
	return dataPoint.ResponseTime
}

func (dataPoint rawHistoricDataPoint) GetTimestamp() int64 {
	return dataPoint.Timestamp
}

func (dataPoint rawHistoricDataPoint) getRawDataPoint() rawHistoricDataPoint {
	return dataPoint
}

func (dataPoint rawHistoricDataPoint) GetStatusCode() int {
	return dataPoint.StatusCode
}
