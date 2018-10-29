package util

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

type modelData struct {
	Pos int
	Tpl []byte
}

type byModelData []modelData

func (a byModelData) Len() int           { return len(a) }
func (a byModelData) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byModelData) Less(i, j int) bool { return a[i].Pos < a[j].Pos }

func UpdateModelFile(name, filename string, fields []*Field) error {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("ParseFile error: %s", err.Error())
	}

	//ast.Print(fset, file)
	//return nil

	v := &modelFileVisitor{file: file, fset: fset, name: name, fields: fields, allFields: fields}
	ast.Walk(v, file)

	v.doModel()
	v.doBModel()
	v.doInitFunc()

	sort.Sort(byModelData(v.modelDatas))

	l := len(v.modelDatas)
	if l == 0 {
		return errors.New("nothing modify")
	}

	b := 0
	e := v.modelDatas[0].Pos
	buf := new(bytes.Buffer)
	for i := 0; i < l; i++ {
		buf.Write(source[b:e])
		buf.Write(v.modelDatas[i].Tpl)
		b = e
		if i+1 < l {
			e = v.modelDatas[i+1].Pos
		}
	}
	buf.Write(source[b:])

	err = ioutil.WriteFile(filename, buf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	Gofmt(filename)
	return nil
}

// struct字段处理
// 增加字段
type modelFileVisitor struct {
	file *ast.File
	fset *token.FileSet

	name      string   // struct Name
	fields    []*Field // 将要新增的字段
	allFields []*Field // 所有的字段（包括本次要新增的）

	mFields      *ast.FieldList // Name struct中的 ast.FieldList
	mFieldsEnd   token.Pos
	bFields      *ast.FieldList // var BName struct中的 ast.FieldList
	bFieldsEnd   token.Pos
	initBVars    []string // init() 函数中的 BName.Field
	initBVarsEnd token.Pos
	modelDatas   []modelData

	importEnd token.Pos
}

func (v *modelFileVisitor) setAllFields(fields []*Field) {
	allFields := []*Field{}
	for _, field := range v.mFields.List {

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

		allFields = append(allFields, f)
	}
	v.allFields = append(allFields, fields...)
}

// Model struct
func (v *modelFileVisitor) doModel() {
	if v.mFieldsEnd == token.NoPos {
		v.mFieldsEnd = v.importEnd
	}

	tpl := []byte{}
	// model 已存在
	if v.mFields != nil {
		lastField := v.mFields.List[v.mFields.NumFields()-1]
		v.mFieldsEnd = lastField.End()
		if lastField.Comment != nil {
			v.mFieldsEnd = lastField.Comment.End()
		}

		fields := distinctFields(v.fields, v.mFields)
		v.setAllFields(fields)
		if len(fields) > 0 {
			tpl = typeStructFieldTpl(fields)
			log.Printf("append new fields to %s model:\n%s\n", v.name, string(tpl))
		} else {
			log.Printf("no new fields to append to %s\n\n", v.name)
		}

	} else {
		tpl = typeStructTpl(v.name, v.fields)
		log.Printf("type %s model:\n%s\n", v.name, string(tpl))
	}

	v.modelDatas = append(v.modelDatas, modelData{Pos: int(v.mFieldsEnd), Tpl: tpl})
}

//var BXxxx struct
func (v *modelFileVisitor) doBModel() {
	if v.bFieldsEnd == token.NoPos {
		v.bFieldsEnd = v.mFieldsEnd + 1
	}

	tpl := []byte{}
	if v.bFields != nil {
		tpl = varBStructFieldTpl(distinctFields(v.allFields, v.bFields))
		log.Printf("append new fields to B%s:\n%s\n", v.name, string(tpl))
	} else {
		tpl = varBStructTpl(v.name, v.allFields)
		log.Printf("type B%s:\n%s\n", v.name, string(tpl))
	}

	v.modelDatas = append(v.modelDatas, modelData{Pos: int(v.bFieldsEnd), Tpl: tpl})
}

// init 函数
func (v *modelFileVisitor) doInitFunc() {
	tpl := []byte{}
	if len(v.initBVars) != 0 || v.initBVarsEnd > 0 {
		tpl = append(tpl, newInitFuncFieldTpl(v.name, v.allFields, v.initBVars)...)
		log.Printf("append new fields to init():\n%s\n", string(tpl))
	} else {

		if v.initBVarsEnd == token.NoPos {
			v.initBVarsEnd = v.mFieldsEnd + 1
		}

		tpl = newInitFuncTpl(v.name, v.allFields, v.initBVars)
		log.Printf("create init() func:\n%s\n", string(tpl))
	}

	v.modelDatas = append(v.modelDatas, modelData{Pos: int(v.initBVarsEnd), Tpl: tpl})
}

func (v *modelFileVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		return v
	}

	switch n := node.(type) {
	case *ast.GenDecl:
		if n.Tok == token.IMPORT {
			v.importEnd = n.End()

		} else if n.Tok == token.TYPE {
			if v.mFieldsEnd == token.NoPos {
				v.mFieldsEnd = n.End()
			}

			ts := n.Specs[0].(*ast.TypeSpec)
			if ts.Name.Name == v.name {
				v.mFields = ts.Type.(*ast.StructType).Fields
				v.mFieldsEnd = ts.End()
			}

		} else if n.Tok == token.VAR {
			ts := n.Specs[0].(*ast.ValueSpec)
			if ts.Names[0].Name == "B"+v.name {
				v.bFieldsEnd = ts.End() - 2
				if s, ok := ts.Type.(*ast.StructType); ok {
					v.bFields = s.Fields
					v.bFieldsEnd = s.Fields.End() - 2
				}
			}
		}

	case *ast.FuncDecl:

		if n.Name.Name == "init" {
			v.initBVarsEnd = n.End() - 2
			for _, l := range n.Body.List {
				if stmt, ok := l.(*ast.AssignStmt); ok {
					for _, lhs := range stmt.Lhs {
						if selExpr, ok := lhs.(*ast.SelectorExpr); ok {
							if x, ok := selExpr.X.(*ast.Ident); ok {
								if x.Name == "B"+v.name {
									v.initBVars = append(v.initBVars, selExpr.Sel.Name)
									v.initBVarsEnd = stmt.End()
								}
							}
						}
					}
				}
			}
		}
	}

	return v
}

func typeStructTpl(name string, fields []*Field) []byte {
	tpl := []byte{}
	tpl = append(tpl, []byte(fmt.Sprintf("\n// %s data struct of %s\ntype %s struct{\n", name, name, name))...)
	tpl = append(tpl, typeStructFieldTpl(fields)...)
	tpl = append(tpl, []byte("}\n")...)
	return tpl
}

func typeStructFieldTpl(fields []*Field) []byte {
	a := []string{}
	for _, field := range fields {
		if field.Comment != "" {
			a = append(a, fmt.Sprintf("    %s %s `bson:\"%s\" json:\"%s\"` // %s\n", field.Name, field.Type, field.TagName, field.TagName, field.Comment))
		} else {
			a = append(a, fmt.Sprintf("    %s %s `bson:\"%s\" json:\"%s\"`\n", field.Name, field.Type, field.TagName, field.TagName))
		}
	}
	return []byte(strings.Join(a, ""))
}

func varBStructTpl(name string, fields []*Field) []byte {
	tpl := []byte{}
	tpl = append(tpl, []byte(fmt.Sprintf("\n// B%s defines %s model bson field name\nvar B%s struct{\n", name, name, name))...)
	tpl = append(tpl, varBStructFieldTpl(fields)...)
	tpl = append(tpl, []byte("}\n")...)
	return tpl
}

func varBStructFieldTpl(fields []*Field) []byte {
	tpl := []byte{}
	for _, f := range fields {
		tpl = append(tpl, []byte(fmt.Sprintf("    %s string\n", f.Name))...)
	}
	return tpl
}

func newInitFuncTpl(name string, fields []*Field, vars []string) []byte {
	tpl := []byte{}
	tpl = append(tpl, []byte("\nfunc init (){\n")...)
	tpl = append(tpl, newInitFuncFieldTpl(name, fields, vars)...)
	tpl = append(tpl, []byte("}\n")...)
	return tpl
}

func newInitFuncFieldTpl(name string, fields []*Field, vars []string) []byte {
	tpl := []byte{}
	var dup bool
	for _, f := range fields {
		dup = false
		for _, v := range vars {
			if f.Name == v {
				dup = true
				break
			}
		}
		if !dup {
			tpl = append(tpl, []byte(fmt.Sprintf("    B%s.%s = \"%s\"\n", name, f.Name, f.TagName))...)
		}
	}
	return tpl
}

func distinctFields(fields []*Field, fieldList *ast.FieldList) []*Field {
	var exists bool
	fs := []*Field{}
	for _, n := range fields {
		exists = false
		for _, f := range fieldList.List {
			if len(f.Names) == 0 {
				continue
			}
			if n.Name == f.Names[0].Name {
				exists = true
				break
			}
		}
		if !exists {
			fs = append(fs, n)
		}
	}
	return fs
}
