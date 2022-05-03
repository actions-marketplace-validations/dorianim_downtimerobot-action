package frontend

import (
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/dorianim/downtimerobot/internal/announcements"
	"github.com/dorianim/downtimerobot/internal/statistics"
	"github.com/dorianim/downtimerobot/internal/templates"
	"github.com/leaanthony/debme"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

//go:embed files/*
var frontendFiles embed.FS
var minifyer = minify.New()

type config struct {
	Frontend frontendConfig
}

type frontendConfig struct {
	Title string
	Icon  string
}

// Generate the static files for the frontend
func Generate(serviceList statistics.ServiceList, serviceDetailList []statistics.ServiceDetails, announcementList *announcements.Announcements) error {
	files, _ := debme.FS(frontendFiles, "files")
	templates, _ := files.FS("templates")
	staticFiles, _ := files.FS("static")
	minifyer.AddFunc("css", css.Minify)
	minifyer.AddFunc("html", html.Minify)
	minifyer.AddFunc("js", js.Minify)

	if err := renderTemplates(templates); err != nil {
		return err
	}
	if err := copyStaticFiles(staticFiles); err != nil {
		return err
	}
	if err := storeServiceList(serviceList); err != nil {
		return err
	}
	if err := storeServiceDetailList(serviceDetailList); err != nil {
		return err
	}
	if err := storeAnnouncementList(announcementList); err != nil {
		return err
	}
	return nil
}

func renderTemplates(templateFiles debme.Debme) error {
	return fs.WalkDir(templateFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".html") || strings.HasPrefix(path, "partials") {
			return nil
		}

		log.WithFields(log.Fields{
			"path": path,
			"name": d.Name(),
		}).Debug("Rendering file")

		parsedTemplate, err := template.ParseFS(templateFiles, path, "partials/*.html")
		if err != nil {
			return err
		}
		config, err := loadConfig()
		if err != nil {
			return err
		}

		result, err := templates.ExecuteTemplate(parsedTemplate, config)
		if err != nil {
			return err
		}
		return minifyAndWriteFile(path, result)
	})
}

func copyStaticFiles(staticFiles debme.Debme) error {
	return fs.WalkDir(staticFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." || d.IsDir() {
			return nil
		}

		log.WithFields(log.Fields{
			"path": path,
			"name": d.Name(),
		}).Debug("Copying file")

		if strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".js") {
			f, err := staticFiles.Open(path)
			if err != nil {
				return err
			}
			b, err := io.ReadAll(f)
			if err != nil {
				return err
			}

			return minifyAndWriteFile("static/"+path, string(b))
		}

		destinationPath := getDestinationPath("static/" + path)
		if err = os.MkdirAll(filepath.Dir(destinationPath), 0777); err != nil {
			return err
		}
		return staticFiles.CopyFile(path, destinationPath, 0666)
	})
}

func minifyAndWriteFile(path string, content string) error {
	pathList := strings.Split(path, ".")
	fileType := pathList[len(pathList)-1]
	content, err := minifyer.String(fileType, content)
	if err != nil {
		return err
	}

	return writeFile(path, []byte(content))
}

func writeFile(path string, content []byte) error {
	destinationPath := getDestinationPath(path)
	if err := os.MkdirAll(filepath.Dir(destinationPath), 0777); err != nil {
		return err
	}

	err := os.WriteFile(destinationPath, content, 0666)
	return err
}

func getDestinationPath(sourcePath string) string {
	return "./public/" + sourcePath
}

func storeServiceList(serviceList statistics.ServiceList) error {
	data, _ := json.MarshalIndent(serviceList, "", " ")
	return writeFile("data/serviceList.json", data)
}

func storeServiceDetailList(serviceDetailList []statistics.ServiceDetails) error {
	for _, serviceDetails := range serviceDetailList {
		data, _ := json.MarshalIndent(serviceDetails, "", " ")
		if err := writeFile("data/"+serviceDetails.Service.Host+".json", data); err != nil {
			return err
		}
	}
	return nil
}

func storeAnnouncementList(announcementList *announcements.Announcements) error {
	data, _ := json.MarshalIndent(announcementList, "", " ")
	return writeFile("data/announcementList.json", data)
}

func loadConfig() (*frontendConfig, error) {
	conf := &config{}
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}
	return &conf.Frontend, nil
}
