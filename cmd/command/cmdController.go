package command

import (
	"fmt"
	"log"
	"os"

	"igen/cmd/util"
)

var cmdController = &CMD{
	UsageLine: "controller [options]",
	Short:     "generate controller",
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
	-ctpl 'default.tpl'
		specify controller tpl file, default is "controller/adm.tpl,controller/vx.tpl" at igen tpl folder
	-hpkg 'name'
		specify router helper package, default is "helper"
	-htpl 'helper.tpl'
		specify router helper package, default is ",controller/helper.tpl" at igen tpl folder

Example:

	igen controller -name=User
	igen controller -name=User -cpkg=adm,v1
	igen controller -name=User -cpkg=v1
	igen controller -name=User -cpkg=v1 -hpkg=helper -htpl=controller/helper.tpl
	`,
	Run: func(a *util.Arg) {
		createControllerFile(a)
		log.Println("done")
	},
}

func init() {
	cmdController.Init()
	CMDs = append(CMDs, cmdController)
}

func createControllerFile(arg *util.Arg) {
	log.Println("---- generate controller ----")

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

	routers := []string{}
	for _, c := range controllers {
		data.CPkgName = c.Pkg
		data.HPkgName = c.HelperPkg

		err := geneControlFile(arg, c, data)
		if err != nil {
			log.Fatalln(err)
			return
		}

		geneHelperFile(arg, c, data)

		routers = append(routers, getRouterStatement(data)...)
	}

	log.Printf("\n\nrouters:\n\n")
	for _, r := range routers {
		fmt.Println(r)
	}
	log.Println("")
}

func getRouterStatement(data *util.TData) []string {
	resource := util.ToUnderScore(data.ModelLowerName)
	return []string{
		fmt.Sprintf("// %s", data.ModelName),
		fmt.Sprintf(`rg.GET("/%ss", %s.List%ss)`, resource, data.CPkgName, data.ModelName),
		fmt.Sprintf(`rg.POST("/%ss", %s.Create%s)`, resource, data.CPkgName, data.ModelName),
		fmt.Sprintf(`rg.GET("/%ss/:id", %s.Show%s)`, resource, data.CPkgName, data.ModelName),
		fmt.Sprintf(`rg.PATCH("/%ss/:id", %s.Update%s)`, resource, data.CPkgName, data.ModelName),
		fmt.Sprintf(`rg.DELETE("/%ss/:id", %s.Delete%s)`, resource, data.CPkgName, data.ModelName),
		"",
	}
}
func geneControlFile(arg *util.Arg, c *util.Controller, data *util.TData) error {
	sourceTpl, destFile, err := c.FilePath()
	if err != nil {
		return err
	}

	if !*arg.Overwrite {
		if _, err := os.Stat(destFile); err == nil {
			return fmt.Errorf("the file %s is exist, if you want to overwrite it, please use -o params\n", destFile)
		}
	}

	log.Println("use control tpl ", sourceTpl)
	if *arg.Overwrite {
		log.Println("overwrite to ", destFile)
	} else {
		log.Println("write to ", destFile)
	}

	err = generateFileFromTPL(data, sourceTpl, destFile)
	if err != nil {
		err = fmt.Errorf("gen controller file error: %s \n", err.Error())
	}

	return err
}

func geneHelperFile(arg *util.Arg, c *util.Controller, data *util.TData) {
	if !c.ValidHelper() {
		log.Println("invalid helper")
		return
	}

	sourceTpl, destFile, err := c.HelperFilePath()
	if err != nil {
		log.Println(err)
		return
	}

	if !*arg.Overwrite {
		if _, err := os.Stat(destFile); err == nil {
			log.Fatalf("the file %s is exist, if you want to overwrite it, please use -o params\n", destFile)
			return
		}
	}

	log.Println("use control helper tpl ", sourceTpl)
	if *arg.Overwrite {
		log.Println("overwrite to ", destFile)
	} else {
		log.Println("write to ", destFile)
	}
	err = generateFileFromTPL(data, sourceTpl, destFile)
	if err != nil {
		log.Fatalf("gen controller helper file error: %s \n", err.Error())
	}
}
