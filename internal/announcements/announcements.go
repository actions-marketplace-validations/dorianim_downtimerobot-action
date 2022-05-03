package announcements

import (
	"time"

	"github.com/spf13/viper"
)

type config struct {
	Announcements announcementsConfig `json:"announcements"`
}

type announcementsConfig struct {
	ExportDays    int               `json:"exportDays"`
	Announcements []rawAnnouncement `json:"announcements"`
}

type rawAnnouncement struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	TimeString string `json:"timeString"`
}

type Announcements struct {
	ExportedDays  int            `json:"exportedDays"`
	Announcements []Announcement `json:"announcements"`
}

type Announcement struct {
	Title      string           `json:"title"`
	Type       AnnouncementType `json:"type"`
	Content    string           `json:"content"`
	TimeString string           `json:"timeString"`
}

func Generate() (*Announcements, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	announcements, err := calculateTimestamps(config.Announcements)
	if err != nil {
		return nil, err
	}

	return announcements, nil
}

func calculateTimestamps(config announcementsConfig) (*Announcements, error) {

	rawAnnouncements := config.Announcements
	announcements := make([]Announcement, 0)
	now := time.Now()

	for i, rawAnnouncement := range rawAnnouncements {
		time, err := time.Parse("2006-01-02 15:04", rawAnnouncement.TimeString)
		if err != nil {
			return nil, err
		}

		if int(now.Sub(time).Hours()) > 24*config.ExportDays {
			continue
		}

		announcements = append(announcements, Announcement{})
		announcements[i].Title = rawAnnouncement.Title
		if err := announcements[i].Type.UnmarshalJSON([]byte("\"" + rawAnnouncement.Type + "\"")); err != nil {
			return nil, err
		}
		announcements[i].Content = rawAnnouncement.Content
		announcements[i].TimeString = time.Format("January 02, 2006, 15:04")
	}

	return &Announcements{Announcements: announcements, ExportedDays: config.ExportDays}, nil
}

func loadConfig() (*config, error) {
	conf := &config{}
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}
	return conf, nil
}
