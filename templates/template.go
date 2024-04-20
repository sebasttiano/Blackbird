// Package templates для работы с html шаблонами.
package templates

import (
	"html/template"
	"log"
	"path/filepath"
	"runtime"
)

// HTMLTemplates хранит *template.Template.
type HTMLTemplates struct {
	IndexTemplate *template.Template
}

// ParseTemplates парсит шаблон html и возвращает HTMLTemplates.
func ParseTemplates() HTMLTemplates {

	ServerTemplates := HTMLTemplates{}
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
