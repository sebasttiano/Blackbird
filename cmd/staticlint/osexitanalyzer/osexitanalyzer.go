package osexitanalyzer

import (
	"golang.org/x/tools/go/analysis"
)

var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitanalyzer",
	Doc:  "check for direct usage of os.Exit in func main package main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	return nil, nil
}
