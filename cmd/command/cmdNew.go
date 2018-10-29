package command

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"igen/cmd/util"
)

const (
	group = "igen"
	demo  = "msdemo"
)

var newArg struct {
	ProjectName *string
	HttpPort    *int
	RpcPort     *int
	ProjectTpl  *string
	Overwrite   *bool
}

var cmdNew = &CMD{
	UsageLine: "new -name=yourProjectName/yourSubProjectName -port=8080 -tpl=msapi",
	Short:     `init a new project`,
	Long: `Options:

	-name 'project name'
		specify the project name, eg: my/foo

	-port 'http port'

	-rpcPort 'rpc port'

	-tpl 'project tpl'
		it has two tpl: msapi, msrpc. msapi is default.
	-o
		overwrite exists files

Example:

	igen new -name=goo/demo -port=8080 -tpl=msapi
	igen new -name=goo/demo -port=8080 -rpcPort=8180 -tpl=msrpc

	`,
	Run: func(a *util.Arg) {
		initProject()
	},
}

func init() {
	newArg.HttpPort = new(int)
	newArg.RpcPort = new(int)
	newArg.ProjectName = new(string)
	newArg.ProjectTpl = new(string)
	newArg.Overwrite = new(bool)

	cmdNew.Flag.StringVar(newArg.ProjectName, "name", "", "")
	cmdNew.Flag.IntVar(newArg.HttpPort, "port", 0, "")
	cmdNew.Flag.IntVar(newArg.RpcPort, "rpcPort", 0, "")
	cmdNew.Flag.StringVar(newArg.ProjectTpl, "tpl", "msapi", "")
	cmdNew.Flag.BoolVar(newArg.Overwrite, "o", false, "")

	CMDs = append(CMDs, cmdNew)
}

func initProject() {
	if newArg.ProjectName == nil || *newArg.ProjectName == "" {
		log.Fatalln("project name can't be empty. please use `igen new -name=yourProjectName/yourSubProjectName`")
		return
	}

	if !strings.Contains(*newArg.ProjectName, "/") {
		log.Fatalln("project name format must be yourProjectName/yourSubProjectName")
		return
	}

	if *newArg.HttpPort == 0 {
		log.Fatalln("http port can't be empty")
		return
	}

	if *newArg.RpcPort == 0 {
		*newArg.RpcPort = *newArg.HttpPort + 100
	}

	log.Printf("init project %s, HTTP port is %d, RPC port is %d\n", *newArg.ProjectName, *newArg.HttpPort, *newArg.RpcPort)
	arr := strings.Split(*newArg.ProjectName, "/")
	err := copyMSDemo(arr[0], arr[1], *newArg.HttpPort, *newArg.RpcPort, *newArg.ProjectTpl, *newArg.Overwrite)
	if err != nil {
		log.Fatalln("fail")
	} else {
		log.Fatalln("done")
	}
}

func copyMSDemo(groupName, projName string, port, rPort int, tpl string, overwrite bool) error {
	log.Println("---- init project ----")

	var newPath string
	var str string

	lowerDemo := strings.ToLower(demo)
	upperDemo := strings.ToUpper(demo)

	lowerProjName := strings.ToLower(projName)
	upperProjName := strings.ToUpper(projName)

	upperGroup := strings.ToUpper(group)
	upperGroupName := strings.ToUpper(groupName)

	demoPath := os.Getenv("GOPATH") + "/src/igen/tpl/" + tpl

	err := filepath.Walk(demoPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		newPathPrefix := os.Getenv("GOPATH") + "/src/" + groupName
		os.Mkdir(newPathPrefix, os.ModePerm)
		newPathPrefix += "/" + projName

		newPath = newPathPrefix + path[len(demoPath):]
		if strings.HasPrefix(newPath, newPathPrefix+"/.git") && newPath != newPathPrefix+"/.gitignore" {
			return nil
		}

		if f.IsDir() {
			log.Printf("[m] %s\n", newPath)
			e := os.Mkdir(newPath, os.ModePerm)
			if e != nil && !overwrite {
				log.Fatalf("mkdir %s error: %s\n", newPath, e.Error())
				return e
			}
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("read file %s error: %s\n", path, err.Error())
			return nil
		}

		if !strings.HasPrefix(newPath, newPathPrefix+"/vendor/") {
			str = string(data)
			str = strings.Replace(str, lowerDemo, lowerProjName, -1)
			str = strings.Replace(str, upperDemo, upperProjName, -1)
			str = strings.Replace(str, group, groupName, -1)
			str = strings.Replace(str, upperGroup, upperGroupName, -1)

			if newPath == newPathPrefix+"/conf/app.toml" ||
				newPath == newPathPrefix+"/Makefile" ||
				newPath == newPathPrefix+"/docker-compose.yml" ||
				newPath == newPathPrefix+"/Dockerfile" {

				str = strings.Replace(str, "8081", strconv.Itoa(port), -1)
				str = strings.Replace(str, "8181", strconv.Itoa(rPort), -1)
			}

			data = []byte(str)
		}

		if newPath == newPathPrefix+"/vendor/vendor.json" {
			str = string(data)
			str = strings.Replace(str, lowerDemo, lowerProjName, -1)
			str = strings.Replace(str, group, groupName, -1)
			data = []byte(str)
		}

		log.Printf("[w] %s\n", newPath)
		err = ioutil.WriteFile(newPath, data, os.ModePerm)
		if err != nil {
			log.Fatalf("write file %s error: %s\n", newPath, err.Error())
			return nil
		}

		return nil
	})

	if err != nil {
		log.Fatalf("filepath.Walk() returned %v\n", err)
		return err
	}

	return nil
}
