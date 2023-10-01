package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	keys := viper.AllKeys()
	println(keys)
	setting := viper.AllSettings()
	fmt.Println(setting)
	server := InitWebServer()

	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})

	server.Run(":8080")
}
