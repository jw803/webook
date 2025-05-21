package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jw803/webook/internal/interface/event"
	"net/http"
	"os"
	"strings"
)

const (
	webapi   = "webapi"
	consumer = "consumer"
)

type App struct {
	webServer map[string]*gin.Engine
	consumers map[string]event.Consumer
}

func (a *App) Start() {
	serviceType, serviceName := a.ExtractServiceInfo()

	switch serviceType {
	case webapi:
		webServer, ok := a.webServer[serviceName]
		if !ok {
			panic("Unsupported AppRole")
		}
		webServer.GET("/health", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "ok")
		})
		webServer.Run(":8081")
	case consumer:
		consumerSvc, ok := a.consumers[serviceName]
		if !ok {
			panic("Unsupported AppRole")
		}
		consumerSvc.Start()
	default:
		panic("Unsupported AppRole")
	}
}

func (a *App) ExtractServiceInfo() (string, string) {
	appRole, ok := os.LookupEnv("APP_ROLE")
	if !ok {
		return "", ""
	}
	appRoleInfo := strings.Split(appRole, "-")
	var serviceType, serviceName string
	serviceType = appRoleInfo[0]
	if len(appRoleInfo) >= 2 {
		serviceType, serviceName = appRoleInfo[0], strings.Join(appRoleInfo[1:], "-")
	}
	return serviceType, serviceName
}
