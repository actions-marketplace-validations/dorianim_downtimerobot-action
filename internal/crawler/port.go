package crawler

import "fmt"

type PortService struct {
	Status      string    `json:"status"`
	Name        string    `json:"name"`
	Url         string    `json:"url"`
	DailyRatios []float32 `json:"dailyRatios"`

	Port int `json:"port"`
}

func (service *PortService) Crawl() error {
	fmt.Println("Port")
	return nil
}
