package templates

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

func ExecuteTemplate(template *template.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	err := template.Funcs(getFuncMap()).Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"join":  strings.Join,
		"print": print,
	}
}

func print(items ...interface{}) string {
	result := ""
	for _, item := range items {
		switch v := item.(type) {
		case string:
			result = fmt.Sprintf("%s%s", result, v)
		case []string:
			result = fmt.Sprintf("%s%s", result, strings.Join(v, ""))
		default:
			continue
		}
	}
	return result
}
