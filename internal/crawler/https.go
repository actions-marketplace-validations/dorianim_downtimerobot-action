package crawler

import (
	"fmt"
	"net/http"
	"time"
)

type httpsService struct {
	genericService `mapstructure:",squash"`

	Path             string `json:"path"`
	ValidStatusCodes []int  `json:"validStatusCodes"`
}

type httpsHistoricDataPoint struct {
	rawHistoricDataPoint
	service *httpsService
}

const httpsRequestError = 600

func (service *httpsService) crawl() HistoricDataPoint {
	var statusCode int
	var statusMessage string
	var responseTime int64

	if service.IsDisabled() {
		statusCode = -1
		statusMessage = "The service is disabled"
		responseTime = -1
	} else {
		start := time.Now()
		resp, err := http.Get(fmt.Sprintf("https://%s%s", service.Host, service.Path))
		if err != nil {
			statusCode = httpsRequestError
			statusMessage = err.Error()
		} else {
			statusCode = resp.StatusCode
			statusMessage = ""
		}
		responseTime = time.Since(start).Milliseconds()
	}

	rawDataPoint := rawHistoricDataPoint{time.Now().Unix(), statusCode, responseTime, statusMessage}
	dataPoint := httpsHistoricDataPoint{rawDataPoint, service}
	service.historicData = append(service.historicData, dataPoint)

	return dataPoint
}

func (service *httpsService) setHistoricData(rawData []rawHistoricDataPoint) {
	service.historicData = make([]HistoricDataPoint, len(rawData))
	for i, rawDataPoint := range rawData {
		service.historicData[i] = httpsHistoricDataPoint{rawDataPoint, service}
	}
}

func (service *httpsService) GetType() string {
	return "https"
}

// == data point ==

func (dataPoint httpsHistoricDataPoint) IsUp() bool {
	if dataPoint.service.ValidStatusCodes == nil {
		return dataPoint.StatusCode >= http.StatusOK && dataPoint.StatusCode <= http.StatusPermanentRedirect
	}
	return sliceContains(dataPoint.service.ValidStatusCodes, dataPoint.StatusCode)
}

func (dataPoint httpsHistoricDataPoint) GetStatusMessage() string {
	message := http.StatusText(dataPoint.StatusCode)
	if len(message) == 0 {
		message = dataPoint.StatusMessage
	}
	return message
}

// == helper ==

func sliceContains(haystack []int, needle int) bool {
	for _, element := range haystack {
		if element == needle {
			return true
		}
	}
	return false
}
