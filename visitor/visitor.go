package visitor

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
	"unicode"
)

const DEBUG = false

type Visitor struct {
	i        int
	FileName string
}

func New(fileName string) Visitor {
	return Visitor{FileName: fileName}
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
	case *ast.TypeSpec:
		if DEBUG {
			fmt.Printf("%s%v\n", strings.Repeat("\t", int(v.i)), t.Name.Name)
			fmt.Printf("%s%T\n", strings.Repeat("\t", int(v.i)), t.Type)
		}
		typeStr := fmt.Sprintf("%T", t.Type)
		switch t2 := t.Type.(type) {
		case *ast.StructType:
			fmt.Println(printStruct(t.Name, t2))
			return nil
		case *ast.ArrayType:
			typeStr = arrayVisit(t2)
		case *ast.Ident:
			typeStr = identVisit(t2)
		case *ast.StarExpr:
			typeStr = starVisit(t2)
		}
		if t.Name.IsExported() {
			fmt.Printf("export type %s = %s\n\n", t.Name.Name, typeStr)
		} else {
			fmt.Printf("type %s = %s\n\n", t.Name.Name, typeStr)
		}
		return nil
	default:
	}

	return Visitor{i: v.i + 1, FileName: v.FileName}
}

func printStruct(name *ast.Ident, st *ast.StructType) string {
	var res string
	if name.IsExported() {
		res = fmt.Sprintf("%s%s", res, fmt.Sprintf("export interface %s {\n", name.Name))
	} else {
		res = fmt.Sprintf("%s%s", res, fmt.Sprintf("interface %s {\n", name.Name))
	}
	res = fmt.Sprintf("%s%s", res, printStructFields(st.Fields))

	res = fmt.Sprintf("%s%s", res, "}\n")
	return res
}

// printStructFields prints the list of fields of a FieldList
func printStructFields(fields *ast.FieldList) string {
	res := ""
	for _, f := range fields.List {
		if f.Names[0].IsExported() {
			var fieldType string
			fieldName := f.Names[0].Name
			if f.Tag != nil {
				tag := f.Tag.Value
				tag = tag[1 : len(tag)-1]
				structTag := reflect.StructTag(tag)
				jsonTag := structTag.Get("json")
				// TODO handle omitempty
				if jsonTag != "" && jsonTag != "-" {
					fieldName = jsonTag
				}
				//fmt.Println(f.Tag.Value)
			}
			switch t := f.Type.(type) {
			case *ast.ArrayType:
				fieldType = arrayVisit(t)
			case *ast.StarExpr:
				fieldType = starVisit(t)
			case *ast.Ident:
				fieldType = identVisit(t)
			default:
				fieldType = fmt.Sprintf("%T", t)
			}
			res = fmt.Sprintf("%s%s", res, fmt.Sprintf("\t%s: %s\n", fieldName, fieldType))
		}
	}
	return res
}

func starVisit(st *ast.StarExpr) string {
	switch t := st.X.(type) {
	case *ast.Ident:
		return identVisit(t)
	default:
		return fmt.Sprintf("%T", t)
	}
}

func identVisit(id *ast.Ident) string {
	switch id.Name {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "number"
	default:
		return id.Name
	}
}

func arrayVisit(arr *ast.ArrayType) string {
	switch t := arr.Elt.(type) {
	case *ast.Ident:
		return fmt.Sprintf("%s[]", identVisit(t))
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
