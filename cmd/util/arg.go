package util

import (
	"errors"
	"fmt"
	"strings"
)

// Arg 创建model, controller 的参数
type Arg struct {
	Name   *string      // 资源名称
	Fields *StringSlice // 字段名

	MTpl     *string // model 的模板
	MComment *string // model 的说明
	MPkgName *string // model 的目录名, 默认是 model
	CRUD     *bool   // 是否同时创建CRUD系列函数

	CDir      *string // 存放controller的目录, 默认是 router
	CPkgNames *string // controller 的目录名, 默认是 adm,v1
	CTpls     *string // controller 的模板
	HPkgName  *string // controller helper package name
	HTpls     *string // controller helper tpl

	Overwrite *bool // weather overwrite the model/controller file

	m *Model
}

func (a *Arg) IsOverwrite() bool {
	if a.Overwrite != nil {
		return *a.Overwrite
	}
	return false
}

func (a *Arg) GetModel() (*Model, error) {
	if a.m != nil {
		return a.m, nil
	}
	var err error
	m := new(Model)

	if a.Name == nil || *a.Name == "" {
		return m, errors.New("name不能为空")
	}

	names := strings.Split(*a.Name, ".")
	nameL := len(names)

	m.Name = strings.ToUpper(names[0][:1]) + names[0][1:]
	m.LowerName = toLowerCamelCase(m.Name)

	m.SName = m.Name
	if nameL == 2 {
		m.SName = strings.ToUpper(names[1][:1]) + names[1][1:]
	}
	m.LowerSName = toLowerCamelCase(m.SName)

	m.Fields, err = a.getFields()
	if err != nil {
		return nil, err
	}

	if a.MPkgName != nil {
		m.Pkg = *a.MPkgName
	}

	if a.MTpl != nil {
		m.Tpl = *a.MTpl
	}

	if a.MComment != nil {
		m.Comment = *a.MComment
	}
	a.m = m
	return a.m, nil
}

func (a *Arg) GetControls() (c []*Controller, err error) {
	if a.CPkgNames == nil {
		err = errors.New("tpl或pkg为空")
		return
	}

	var model *Model
	model, err = a.GetModel()
	if err != nil {
		return
	}

	pkgNames := strings.Split(*a.CPkgNames, ",")

	tpls := []string{}
	if a.CTpls != nil {
		tpls = strings.Split(*a.CTpls, ",")
	}
	tl := len(tpls)

	helperTpls := []string{}
	if a.HTpls != nil {
		helperTpls = strings.Split(*a.HTpls, ",")
	}
	hl := len(helperTpls)

	for i, pkg := range pkgNames {
		ctrl := &Controller{}
		ctrl.Dir = *a.CDir
		ctrl.Pkg = pkg
		ctrl.Model = model

		if i < tl && tpls[i] != "" {
			ctrl.Tpl = tpls[i]
		}

		// helper
		if a.HPkgName != nil {
			ctrl.HelperPkg = *a.HPkgName
		}
		if i < hl && helperTpls[i] != "" {
			ctrl.HelperTpl = helperTpls[i]
		}

		c = append(c, ctrl)
	}

	return
}

func (a *Arg) getFields() ([]*Field, error) {
	fields := []*Field{}
	for _, f := range *a.Fields {
		arr := strings.Split(f, ":")
		l := len(arr)
		if l < 2 {
			return fields, fmt.Errorf("%s field format error\n", f)
		}

		key := arr[0]
		key = strings.ToUpper(key[:1]) + key[1:]

		field := &Field{
			Name:    key,
			Type:    arr[1],
			TagName: toLowerCamelCase(key),
		}
		if l > 2 {
			field.Comment = arr[2]
		}

		fields = append(fields, field)
	}
	return fields, nil
}

func (a *Arg) ToTData() (*TData, error) {

	model, err := a.GetModel()
	if err != nil {
		return nil, err
	}

	data := new(TData)
	data.ProjectName = getProjectName()
	data.SubProjectName = getSubProjectName()
	data.Fields = model.Fields
	data.ModelName = model.SName
	data.ModelLowerName = model.LowerSName
	data.ModelComment = model.Comment
	data.MPkgName = model.Pkg
	data.CollectionName = fmt.Sprintf("%s.%ss", getSubProjectName(), strings.ToLower(data.ModelName))
	if data.ModelComment == "" {
		data.ModelComment = fmt.Sprintf("defines %s's data structure", data.ModelName)
	}
	if a.CDir != nil {
		data.CDir = *a.CDir
	}
	if a.CRUD != nil {
		data.CRUD = *a.CRUD
	}

	return data, nil
}
