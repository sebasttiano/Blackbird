// Package osexitanalyzer содержит кастомный линтер, проверяющий использование ф-ции os.Exit.
package osexitanalyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer переменная типа Analyzer
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexitanalyzer",
	Doc:  "check for direct usage of os.Exit in func main package main",
	Run:  run,
}

// run запуск анализатора
func run(pass *analysis.Pass) (interface{}, error) {

	checkExit := func(f *ast.FuncDecl) {
		for _, stmt := range f.Body.List {
			stmt, ok := stmt.(*ast.ExprStmt)
			if !ok {
				continue
			}
			call, ok := stmt.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			if selector.Sel.Name == "Exit" {
				pass.Reportf(selector.Sel.NamePos, "call of exit func in main")
			}
		}
	}
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			p, ok := node.(*ast.File)
			if ok && p.Name.Name == "main" {
				ast.Inspect(p, func(innerNode ast.Node) bool {
					if x, ok := innerNode.(*ast.FuncDecl); ok {
						if x.Name.Name == "main" {
							checkExit(x)
						}
					}
					return true
				})
			}
			return true
		})
	}
	return nil, nil
}
