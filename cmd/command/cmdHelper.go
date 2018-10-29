package command

import (
	"fmt"
	"log"

	"igen/cmd/util"
)

var cmdHelper = &CMD{
	UsageLine: "helper [command]",
	Short:     "generate controller helper file",
	Run: func(a *util.Arg) {
		createHelperFile(a)
	},
	Long: `Options:

	-name 'model name'
		specify the model name, eg: Foo
	-f 'field list'
		list of fields, the format is fieldName:fieldType:comment
		eg: -f name:string -f age:18:Age
	-o
		overwrite the model file if it is exists

	-mpkg 'name'
		specify model package, default is "model"

	-cpkg 'package name'
		specify controller sub package, default is "adm,v1"
	-cdir 'controller package name'
		specify controllers package, default is "router"
	-hpkg 'name'
		specify router helper package, default is "helper"
	-htpl 'helper.tpl'
		specify router helper package, default is ",controller/helper.tpl" at igen tpl folder

Example:

	igen helper -name=User -cpkg=adm
	igen helper -name=User -cpkg=adm,v1
	`,
}

func init() {
	cmdHelper.Init()
	CMDs = append(CMDs, cmdHelper)
}

func createHelperFile(arg *util.Arg) {
	log.Println("---- generate controller helper ----")

	model, err := arg.GetModel()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	controllers, err := arg.GetControls()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	if len(controllers) == 0 {
		log.Fatalln("not found controller package")
		return
	}

	data, err := arg.ToTData()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	if len(data.Fields) == 0 {
		data.Fields, err = model.ParseFields()
		if err != nil {
			log.Fatalf("parse fields error: %s", err.Error())
			return
		}
	}

	for _, c := range controllers {
		data.CPkgName = c.Pkg
		data.HPkgName = c.HelperPkg
		if c.HelperTpl == "" {
			c.HelperTpl = defaultHelperTpl
		}
		geneHelperFile(arg, c, data)
	}

	fmt.Println("done")
}
