// Package main multichecker собирает все анализаторы для совместного запуска
package main

import (
	"encoding/json"
	gocritic "github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/sebasttiano/Blackbird.git/cmd/staticlint/osexitanalyzer"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"os"
	"path/filepath"
)

// Config — имя файла конфигурации.
const Config = `config.json`

// ConfigData описывает структуру файла конфигурации.
type ConfigData struct {
	SimpleCheck []string `json:"simple"`
	StyleCheck  []string
	QuickFix    []string
}

func main() {
	appfile, err := os.Executable()
	if err != nil {
		logger.Log.Error("failed to get executable path", zap.Error(err))
		return
	}

	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), Config))
	if err != nil {
		logger.Log.Error("failed to read config file", zap.Error(err))
	}

	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		logger.Log.Error("failed to parse config file", zap.Error(err))
	}

	mychecks := []*analysis.Analyzer{
		osexitanalyzer.OsExitAnalyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		ineffassign.Analyzer,
		gocritic.Analyzer,
	}

	checks := make(map[string]bool)
	for _, v := range cfg.SimpleCheck {
		checks[v] = true
	}
	for _, v := range cfg.StyleCheck {
		checks[v] = true
	}
	for _, v := range cfg.QuickFix {
		checks[v] = true
	}

	// добавляем все анализаторы из staticcheck simple, которые указаны в файле конфигурации
	for _, v := range simple.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	// добавляем все анализаторы из staticcheck stylecheck, которые указаны в файле конфигурации
	for _, v := range stylecheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	// добавляем все анализаторы из staticcheck quickfix, которые указаны в файле конфигурации
	for _, v := range quickfix.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	// добавляем все анализаторы из staticcheck
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	multichecker.Main(
		mychecks...,
	)
}
