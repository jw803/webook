package main

import (
	"github.com/jw803/webook/config"
)

func main() {
	config.Init()
	app := InitApp()
	app.Start()
}
