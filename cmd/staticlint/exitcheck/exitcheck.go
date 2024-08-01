package exitcheck

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for use os.Exit() in main() func",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {

			if file.Name.Name != "main" {
				return true
			}

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
				if f.Name.Name != "main" {
					return true
				}
				for _, stmt := range f.Body.List {
					if node, ok := stmt.(*ast.ExprStmt); ok {
						if f, ok := node.X.(*ast.CallExpr); ok {
							if ident, ok := f.Fun.(*ast.SelectorExpr); ok {
								if i, ok := ident.X.(*ast.Ident); ok && i.Name == "os" && ident.Sel.Name == "Exit" {
									pass.Reportf(f.Pos(), "calling os.Exit() in main() function")
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
