package command

import (
	"flag"
	"os"
	"path"
	"text/template"

	"igen/cmd/util"
)

const (
	defaultModelPackage   string = "model"
	defaultModelTpl       string = "model/default.tpl"
	defaultControllerDir  string = "router"
	defaultControllerPkgs string = "adm,v1"
	defaultControllerTpls string = "controller/adm.tpl,controller/vx.tpl"
	defaultHelperPkg      string = "helper"
	defaultHelperTpl      string = "controller/helper.tpl"
	defaultHelperTpls     string = ",controller/helper.tpl"
)

func initArg() (*util.Arg, flag.FlagSet) {
	a := new(util.Arg)
	a.Fields = new(util.StringSlice)
	a.Name = new(string)

	a.Overwrite = new(bool)

	a.MTpl = new(string)
	a.MPkgName = new(string)
	a.MComment = new(string)

	a.CDir = new(string)
	a.CPkgNames = new(string)
	a.CTpls = new(string)
	a.HPkgName = new(string)
	a.HTpls = new(string)

	a.CRUD = new(bool)

	var f flag.FlagSet

	f.StringVar(a.Name, "name", "", "")
	f.Var(a.Fields, "f", "")

	f.BoolVar(a.Overwrite, "o", false, "")

	f.StringVar(a.MPkgName, "mpkg", defaultModelPackage, "")
	f.StringVar(a.MTpl, "mtpl", defaultModelTpl, "")
	f.StringVar(a.MComment, "mc", "", "")

	f.StringVar(a.CDir, "cdir", defaultControllerDir, "")
	f.StringVar(a.CPkgNames, "cpkg", defaultControllerPkgs, "")
	f.StringVar(a.CTpls, "ctpl", defaultControllerTpls, "")
	f.StringVar(a.HPkgName, "hpkg", defaultHelperPkg, "")
	f.StringVar(a.HTpls, "htpl", defaultHelperTpls, "")

	f.BoolVar(a.CRUD, "crud", true, "")

	return a, f
}

func generateFileFromTPL(data interface{}, tplPath, dest string) error {
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return err
	}

	file, err := ensureOpenFile(dest)
	if err != nil {
		return err
	}

	if err = tpl.Execute(file, data); err != nil {
		return err
	}

	//format the code
	util.Gofmt(dest)

	return nil
}

func ensureOpenFile(filePath string) (*os.File, error) {
	err := os.MkdirAll(path.Dir(filePath), 0777)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0777)
}
