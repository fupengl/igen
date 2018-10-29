package command

import (
	"log"
	"os"

	"igen/cmd/util"
)

var cmdModel = &CMD{
	UsageLine: "model [options]",
	Short:     "generate model",
	Long: `Options:

	-name 'model name'
		specify the model name, eg: Foo
	-f 'field list'
		list of fields, the format is fieldName:fieldType:comment
		eg: -f name:string -f age:18:Age
	-o
		overwrite the model file if it is exists
	-mc 'model comment'
		specify model comment
	-mtpl 'default.tpl'
		specify model tpl file, default is "model/default.tpl"" at igen tpl folder
	-mpkg 'name'
		specify model package, default is "model"
	-crud
		weather create CRUD functions

Example:

	igen model -name=User -f=username:string:姓名 -f=age:int:年龄 -f=password:string
	igen model -name=User -f=gender:int:性别
	igen model -name=User.Address -f=cityName:string:城市名称 -f=cityId:int

	`,
	Run: func(a *util.Arg) {
		implModel(a)
		log.Println("done")
	},
}

func init() {
	cmdModel.Init()
	CMDs = append(CMDs, cmdModel)
}

func implModel(arg *util.Arg) {
	model, err := arg.GetModel()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	destFile := model.SourcePath()

	fileExists := false
	if _, err := os.Stat(destFile); err == nil {
		fileExists = true
	}

	if model.Name != model.SName {
		if fileExists {
			// 添加或更改struct
			err = util.UpdateModelFile(model.SName, model.SourcePath(), model.Fields)
			if err != nil {
				log.Fatalln(err)
			}

		} else {
			log.Fatalf("the file %s isn't exist\n", destFile)
			return
		}

	} else {

		if fileExists && !arg.IsOverwrite() {
			// 更改 model
			err = util.UpdateModelFile(model.SName, model.SourcePath(), model.Fields)
			if err != nil {
				log.Fatalln(err)
			}

		} else {
			log.Println("---- generate model ----")
			if len(model.Fields) == 0 {
				log.Fatalln("fields can't be empty")
				return
			}
			createModelFile(model, arg)
		}
	}
}

func createModelFile(model *util.Model, arg *util.Arg) {
	data, err := arg.ToTData()
	if err != nil {
		log.Fatalln(err)
		return
	}

	modelTpl, destFile, err := model.FilePath()
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println("use model tpl ", modelTpl)
	if arg.IsOverwrite() {
		log.Println("overwrite to ", destFile)
	} else {
		log.Println("write to ", destFile)
	}

	err = generateFileFromTPL(data, modelTpl, destFile)
	if err != nil {
		log.Fatalf("gen model file error: %s \n", err.Error())
	}
}
