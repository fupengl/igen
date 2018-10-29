package command

import (
	"log"

	"igen/cmd/util"
)

var cmdAll = &CMD{
	UsageLine: "all [options]",
	Short:     "generate both model and controller",
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
	-mtpl 'default.tpl'
		specify model tpl file, default is "model/default.tpl"" at igen tpl folder
	-mc 'model comment'
		specify model comment

	-cpkg 'package name'
		specify controller sub package, default is "adm" and the "v1"
	-cdir 'controller package name'
		specify controllers package, default is "router"
	-ctpl 'default.tpl'
		specify controller tpl file, default is "controller/adm.tpl,controller/vx.tpl" at igen tpl folder
	-hpkg 'name'
		specify router helper package, default is "helper"
	-htpl 'helper.tpl'
		specify router helper package, default is ",controller/helper.tpl" at igen tpl folder

Example:

	igen all -name=User -f=username:string:姓名 -f=age:int:年龄 -f=password:string
	igen all -name=User -f=username:string:姓名 -f=age:int:年龄 -f=password:string -cpkg=adm

	`,

	Run: func(a *util.Arg) {
		implModel(a)
		createControllerFile(a)
		log.Println("done")
	},
}

func init() {
	cmdAll.Init()
	CMDs = append(CMDs, cmdAll)
}
