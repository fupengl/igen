package util

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

func GetStructFields(name, filename string) ([]*Field, error) {
	out := []*Field{}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return out, fmt.Errorf("ParseFile error: %s", err.Error())
	}

	v := &structFieldVisitor{file: file, fset: fset, name: name}
	v.handle = func(fl *ast.FieldList) {
		for _, field := range fl.List {

			f := new(Field)
			if len(field.Names) > 0 {
				f.Name = field.Names[0].Name
			}

			// skip third field
			if f.Name == "" || f.Name == "ID" || f.Name == "UpdatedAt" || f.Name == "CreatedAt" {
				continue
			}

			f.TagName = toLowerCamelCase(f.Name)
			f.Comment = strings.Trim(field.Comment.Text(), "\n")

			var buf bytes.Buffer
			printer.Fprint(&buf, fset, field.Type)
			f.Type = buf.String()

			out = append(out, f)
		}
	}

	ast.Walk(v, file)

	return out, err
}

// struct字段处理
// 增加字段
type structFieldVisitor struct {
	file *ast.File
	fset *token.FileSet

	name   string // struct name
	handle func(fields *ast.FieldList)
}

func (v *structFieldVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		return v
	}

	switch n := node.(type) {
	case *ast.GenDecl:
		if n.Tok != token.TYPE {
			break
		}
		ts := n.Specs[0].(*ast.TypeSpec)
		if ts.Name.Name == v.name {
			fields := ts.Type.(*ast.StructType).Fields
			v.handle(fields)
		}
	}

	return v
}
