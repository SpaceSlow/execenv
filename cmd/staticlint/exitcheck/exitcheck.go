package exitcheck

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer Анализатор для предупреждения прямого вызова os.Exit() в функции main() модуля main.
var Analyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for use os.Exit() in main() func",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {

			// Проверка принадлежности к package main
			if file.Name.Name != "main" {
				return true
			}

			// Проверка наличия импорта стандартной библиотеки os
			hasOSImport := false
			for _, i := range file.Imports {
				if i.Path.Value == `"os"` {
					hasOSImport = true
					break
				}
			}

			if !hasOSImport {
				return true
			}

			if f, ok := n.(*ast.FuncDecl); ok {
				// Проверка наличия функции main
				if f.Name.Name != "main" {
					return true
				}
				for _, stmt := range f.Body.List {
					if node, ok := stmt.(*ast.ExprStmt); ok {
						if f, ok := node.X.(*ast.CallExpr); ok {
							if ident, ok := f.Fun.(*ast.SelectorExpr); ok {
								// Поиск в блоке функции main запуск os.Exit() функции
								if i, ok := ident.X.(*ast.Ident); ok && i.Name == "os" && ident.Sel.Name == "Exit" {
									pass.Reportf(i.NamePos, "calling os.Exit func in main func")
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
