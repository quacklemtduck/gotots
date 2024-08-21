package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"github.com/quacklemtduck/gotots/visitor"
)

func main() {

	fmt.Printf("Generating for file %s\n", os.Getenv("GOFILE"))
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	file := fmt.Sprintf("%s/%s", cwd, os.Getenv("GOFILE"))
	fmt.Printf("Full path: %s\n", file)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(f.Name.Name)
	fmt.Printf("%#v\n", f)
	var v visitor.Visitor
	ast.Walk(v, f)

}
