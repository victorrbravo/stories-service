package tools

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func CollectTypes(path string) map[string]ast.Expr {
	var err error
	var packages map[string]*ast.Package
	log.Println(os.Getwd())
	packages, err = parser.ParseDir(&token.FileSet{}, path, func(fi os.FileInfo) bool {
		log.Println(fi.Name())
		return true
	}, parser.ParseComments)
	if err != nil {
		return make(map[string]ast.Expr)
	}
	log.Println(packages)
	return make(map[string]ast.Expr)
}
