package crawler

import "fmt"

type patternService struct {
	Status      string    `json:"status"`
	Name        string    `json:"name"`
	DailyRatios []float32 `json:"dailyRatios"`

	Path    string `json:"path"`
	Pattern string `json:"pattern"`
}

func (service *patternService) Crawl() error {
	fmt.Println("Pattern")
	return nil
}
