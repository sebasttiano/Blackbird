package templates

import (
	"html/template"
	"log"
	"path/filepath"
	"runtime"
)

type HtmlTemplates struct {
	IndexTemplate *template.Template
}

func ParseTemplates() HtmlTemplates {

	ServerTemplates := HtmlTemplates{}
	var err error
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(b)
	)

	ServerTemplates.IndexTemplate, err = template.ParseFiles(basepath + "/index.html")
	if err != nil {
		log.Fatalf("Couldn`t parse templates %v", err)
	}
	return ServerTemplates
}
