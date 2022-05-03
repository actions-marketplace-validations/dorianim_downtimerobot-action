package notifications

import (
	"errors"
	"regexp"
	"text/template"

	"github.com/dorianim/downtimerobot/internal/crawler"
	"github.com/dorianim/downtimerobot/internal/templates"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type notificationConfig struct {
	NotificationTargets []notificationTarget `json:"notificationTargets"`
}

type notificationTarget struct {
	Enabled         bool   `json:"enabled"`
	Name            string `json:"name"`
	Template        string `json:"template"`
	ShoutrrrURL     string `json:"shoutrrrUrl"`
	ServicesPattern string `json:"servicesPattern"`
}

var config *notificationConfig = nil

func Notify(crawledServices []crawler.Service) error {
	var err error
	config, err = loadConfig()
	if err != nil {
		return err
	}

	return notifyServices(crawledServices)
}

func notifyServices(services []crawler.Service) error {
	for _, service := range services {
		if err := notifyService(service); err != nil {
			log.WithFields(log.Fields{
				"service": service.GetHost(),
				"err":     err.Error(),
			}).Error("Error sending notification to service")
		}
	}
	return nil
}

func notifyService(service crawler.Service) error {
	historicData := service.GetHistoricData()
	historicDataLength := len(historicData)

	if len(historicData) >= 2 && historicData[historicDataLength-1].IsUp() != historicData[historicDataLength-2].IsUp() {
		log.WithFields(log.Fields{
			"service": service,
		}).Debug("Service changed state")
		return sendNotificationForService(service)
	}

	return nil
}

func sendNotificationForService(service crawler.Service) error {
	targets, err := getNotificationTargetsForService(service)
	if err != nil {
		return err
	}

	for _, target := range targets {
		if err := sendNotificationForServiceToTarget(service, target); err != nil {
			log.WithFields(log.Fields{
				"service":            service.GetHost(),
				"notificationTarget": target.Name,
				"err":                err.Error(),
			}).Error("Error sending notification to target")
		}
	}

	return nil
}

func sendNotificationForServiceToTarget(service crawler.Service, target notificationTarget) error {
	parsedTemplate, err := template.New("t").Parse(target.Template)
	if err != nil {
		log.WithFields(log.Fields{
			"service":            service.GetHost(),
			"notificationTarget": target.Name,
			"template":           target.Template,
			"err":                err.Error(),
		}).Error("Error parsing template")
		return errors.New("Error parsing template")
	}

	service.IsUp()
	executedTemplate, err := templates.ExecuteTemplate(parsedTemplate, map[string]interface{}{
		"Service": service,
		"Target":  target,
	})

	//return shoutrrr.Send(target.ShoutrrrURL, executedTemplate)

	log.Info("Sending: " + executedTemplate)
	return nil
}

func getNotificationTargetsForService(service crawler.Service) ([]notificationTarget, error) {
	targets := make([]notificationTarget, 0)
	for _, target := range config.NotificationTargets {
		matched, err := regexp.Match(target.ServicesPattern, []byte(service.GetName()))
		if err != nil {
			log.WithFields(log.Fields{
				"service":            service.GetName(),
				"notificationTarget": target.Name,
				"servicesPattern":    target.ServicesPattern,
				"err":                err,
			}).Error("Failed to compile service pattern")
		} else if matched {
			targets = append(targets, target)
		}
	}
	return targets, nil
}

func loadConfig() (*notificationConfig, error) {
	conf := &notificationConfig{}
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}
	return conf, nil
}
