package visitor

import (
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strings"
	"unicode"
)

const DEBUG = false

type Visitor struct {
	i           int
	FileName    string
	UseComments bool
}

func New(fileName string, useComments bool) Visitor {
	return Visitor{FileName: fileName, UseComments: useComments}
}

func (v Visitor) Visit(n ast.Node) ast.Visitor {

	if n == nil {
		return nil
	}

	if DEBUG {
		fmt.Printf("%s%T\n", strings.Repeat("\t", int(v.i)), n)
	}
	switch t := n.(type) {
	case *ast.File:
		fmt.Printf("// Generated from file: %s\n", v.FileName)
		fmt.Printf("// Package name: %s\n\n", t.Name.Name)
	case *ast.GenDecl:
		if v.UseComments && t.Tok == token.TYPE && t.Doc != nil {
			fmt.Print(v.visitCommentGroup(t.Doc, ""))
		}
	case *ast.TypeSpec:
		if DEBUG {
			fmt.Printf("%s%v\n", strings.Repeat("\t", int(v.i)), t.Name.Name)
			fmt.Printf("%s%T\n", strings.Repeat("\t", int(v.i)), t.Type)
		}
		var typeStr string
		switch t2 := t.Type.(type) {
		case *ast.InterfaceType:
			return nil
		default:
			typeStr = v.visitExpr(t2)
		}
		if t.Name.IsExported() {
			fmt.Printf("export type %s = %s\n\n", t.Name.Name, typeStr)
		} else {
			fmt.Printf("type %s = %s\n\n", t.Name.Name, typeStr)
		}
		return nil
	default:
		if DEBUG {
			fmt.Printf("%s%T Default\n", strings.Repeat("\t", int(v.i)), t)
			fmt.Printf("%s%#v Default\n", strings.Repeat("\t", int(v.i)), t)
		}
	}

	return Visitor{i: v.i + 1, FileName: v.FileName, UseComments: v.UseComments}
}

func (v Visitor) visitExpr(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return v.identVisit(t)
	case *ast.ArrayType:
		return v.arrayVisit(t)
	case *ast.MapType:
		return v.mapVisit(t)
	case *ast.StarExpr:
		return v.starVisit(t)
	case *ast.StructType:
		return v.visitStruct(t)
	}

	return fmt.Sprintf("%T", e)
}

func (v Visitor) visitStruct(st *ast.StructType) string {
	var res string
	res = "{\n"
	res = fmt.Sprintf("%s%s", res, v.printStructFields(st.Fields))
	res = fmt.Sprintf("%s%s", res, "}")
	return res
}

// printStructFields prints the list of fields of a FieldList
func (v Visitor) printStructFields(fields *ast.FieldList) string {
	res := ""
	for _, f := range fields.List {
		for i := range f.Names {
			if f.Names[i].IsExported() {
				var fieldType string
				fieldName := f.Names[i].Name
				extra := ""
				if f.Tag != nil {
					tag := f.Tag.Value
					tag = tag[1 : len(tag)-1]
					structTag := reflect.StructTag(tag)
					jsonTag := structTag.Get("json")

					if jsonTag != "" {
						jsonValues := strings.Split(jsonTag, ",")
						if jsonValues[0] == "-" && len(jsonValues) == 1 {
							continue
						}
						if jsonValues[0] != "" {
							fieldName = jsonValues[0]
						}
						if len(jsonValues) >= 2 && jsonValues[1] == "omitempty" {
							extra = "?"
						}
					}
					//fmt.Println(f.Tag.Value)
				}
				if v.UseComments && f.Doc != nil {
					res = fmt.Sprintf("%s%s", res, v.visitCommentGroup(f.Doc, "\t"))
				}
				fieldType = v.visitExpr(f.Type)
				res = fmt.Sprintf("%s%s", res, fmt.Sprintf("\t%s%s: %s\n", fieldName, extra, fieldType))
			}
		}
	}
	return res
}

func (v Visitor) mapVisit(m *ast.MapType) string {
	return fmt.Sprintf("Record<%s, %s>", v.visitExpr(m.Key), v.visitExpr(m.Value))
}

func (v Visitor) starVisit(st *ast.StarExpr) string {
	switch t := st.X.(type) {
	case *ast.Ident:
		return fmt.Sprintf("%s | null", v.identVisit(t))
	default:
		return fmt.Sprintf("%T", t)
	}
}

func (v Visitor) identVisit(id *ast.Ident) string {
	switch id.Name {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "byte":
		return "number"
	case "bool":
		return "boolean"
	case "rune":
		return "number"
	default:
		return id.Name
	}
}

func (v Visitor) arrayVisit(arr *ast.ArrayType) string {
	switch t := arr.Elt.(type) {
	case *ast.Ident:
		return fmt.Sprintf("%s[]", v.identVisit(t))
	default:
		return fmt.Sprintf("%T[]", arr.Elt)
	}
}

// IsFirstLetterCapitalized returns true if the first letter is capitalized
func IsFirstLetterCapitalized(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsUpper(rune(s[0]))
}

// DecapitalizeFirstLetter returns the string with the first letter decapitalized
func DecapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}

	// Get the first rune (character)
	firstRune := rune(s[0])

	// If it's already lowercase, return the original string
	if unicode.IsLower(firstRune) {
		return s
	}

	// Decapitalize the first letter and concatenate with the rest of the string
	return string(unicode.ToLower(firstRune)) + s[1:]
}

func (v Visitor) visitCommentGroup(doc *ast.CommentGroup, prefix string) string {
	res := fmt.Sprintf("%s/**\n", prefix)
	lines := strings.Split(doc.Text(), "\n")
	for _, l := range lines[:len(lines)-1] {
		res = fmt.Sprintf("%s%s * %s\n", res, prefix, l)
	}
	res = fmt.Sprintf("%s%s */\n", res, prefix)
	return res
}
