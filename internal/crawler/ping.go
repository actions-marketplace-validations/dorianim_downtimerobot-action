package crawler

import "fmt"

type PingService struct {
	Status      string    `json:"status"`
	Name        string    `json:"name"`
	Url         string    `json:"url"`
	DailyRatios []float32 `json:"dailyRatios"`
}

func (service *PingService) Crawl() error {
	fmt.Println("Ping")

	return nil
}
