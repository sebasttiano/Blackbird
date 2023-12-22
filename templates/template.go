package templates

import (
	"html/template"
	"log"
)

type HtmlTemplates struct {
	IndexTemplate *template.Template
}

func ParseTemplates() HtmlTemplates {

	ServerTemplates := HtmlTemplates{}
	var err error
	ServerTemplates.IndexTemplate, err = template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("Couldn`t parse templates %v", err)
	}
	return ServerTemplates
}
