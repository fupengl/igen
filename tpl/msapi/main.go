package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"

	"igen/lib"
	"igen/lib/constant"
	"igen/lib/consul"
	"igen/lib/logger"
	"igen/msdemo/conf"
	"igen/msdemo/router"
)

func main() {
	conf.Init()

	handleSignal()
	servingHTTP()
}

func handleSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Println("Exit ", s)
				friendlyExit()
			case syscall.SIGUSR1:
				log.Println("USR1", s)
			case syscall.SIGUSR2:
				log.Println("USR2", s)
			default:
				log.Println("OTHER Exist", s)
			}
		}
	}()
}

func friendlyExit() {
	consul.DeregisterAll()

	os.Exit(0)
}

// HTTP server
func servingHTTP() {
	r := gin.New()

	if conf.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(lib.GinLogger(nil))
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, POST, PATCH, PUT, DELETE",
		RequestHeaders:  fmt.Sprintf("Origin, Authorization, Content-Type, x-access-token, %s, %s, %s", constant.SignKey, constant.SignVal, constant.SignTime),
		ExposedHeaders:  "",
		MaxAge:          30 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	// Set routers
	router.Init(r)

	// Run
	addr := fmt.Sprintf("%s:%d", conf.App.HTTPAddress, conf.App.HTTPPort)
	logger.Infof("listen HTTP %s", addr)

	r.Run(addr)
}
