package visitor

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"
)

type Visitor struct {
	name string
	i    int
}

func (v Visitor) Visit(n ast.Node) ast.Visitor {

	if n == nil {
		return nil
	}

	fmt.Printf("%s%T\n", strings.Repeat("\t", int(v.i)), n)
	switch t := n.(type) {
	case *ast.Ident:
		fmt.Printf("%s%v\n", strings.Repeat("\t", int(v.i)), t.Name)
		return Visitor{i: v.i + 1, name: t.Name}
	case *ast.TypeSpec:
		fmt.Printf("%s%v\n", strings.Repeat("\t", int(v.i)), t.Name.Name)
		return Visitor{i: v.i + 1, name: t.Name.Name}
	case *ast.StructType:
		fmt.Printf("%s%v\n", strings.Repeat("\t", int(v.i)), len(t.Fields.List))
		fmt.Printf("interface %s {\n", v.name)
		printFields(t.Fields)
		fmt.Printf("}\n")
		return nil
	default:
	}

	return Visitor{i: v.i + 1, name: v.name}
}

// printFields prints the list of fields of a FieldList
func printFields(fields *ast.FieldList) {
	for _, f := range fields.List {
		if f.Names[0].IsExported() {
			var fieldType string
			switch t := f.Type.(type) {
			case *ast.ArrayType:
				fieldType = fmt.Sprintf("%s[]", t.Elt)
			default:
				fieldType = fmt.Sprintf("%s", t)
			}
			fmt.Printf("\t%v: %s", f.Names[0].Name, fieldType)
			fmt.Printf("\n")
		}
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
