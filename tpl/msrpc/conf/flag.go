package conf

import (
	"flag"
	"fmt"
	"os"
)

var clParam *commandLineParam

// 命令行参数
type commandLineParam struct {
	Env         *string
	LoggerLevel *string
	LogFilename *string
	Consul      *string
	Help        *bool
}

func init() {
	clParam = &commandLineParam{
		Env:         flag.String("env", "", "运行环境: prod, dev"),
		LoggerLevel: flag.String("logLevel", "", "日志打印级别: debug, warn, error"),
		LogFilename: flag.String("logFile", "", "日志目录及文件名(默认为 log/access.log)"),
		Consul:      flag.String("consul", "", "consul的地址127.0.0.1:8500"),
		Help:        flag.Bool("help", false, "help"),
	}
}

func loadFlagConf() {
	if clParam == nil {
		return
	}

	if *clParam.Env != "" {
		if *clParam.Env == EnvDev || *clParam.Env == EnvProd || *clParam.Env == EnvTest {
			App.Env = *clParam.Env
		} else {
			fmt.Print("运行环境设置不对: prod, dev\n")
			os.Exit(0)
		}
	}

	if *clParam.LoggerLevel != "" {
		App.LogLevel = *clParam.LoggerLevel
	}

	if *clParam.LogFilename != "" {
		App.LogFilename = *clParam.LogFilename

		App.useFileLogger = true
	}

	if *clParam.Consul != "" {
		App.Consul.Address = *clParam.Consul
	}
}

func consulAddress() string {
	if *clParam.Consul != "" {
		return *clParam.Consul
	}
	if v := os.Getenv("IGEN_CONSUL_ADDRESS"); v != "" {
		return v
	}
	return ""
}

// 命令行参数
func parseFlag() {
	flag.Usage = func() {
		fmt.Print("\n命令行参数设置\n")
		fmt.Fprintf(os.Stdout, "Usage of %s [options] \n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("")
		os.Exit(0)
	}

	if !flag.Parsed() {
		flag.Parse()
	}

	if *clParam.Help {
		flag.Usage()
		os.Exit(0)
	}
}
