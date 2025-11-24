package main

import (
	"github.com/kgugunava/avito-tech-internship/internal/app"
)


func main() {
	app := app.NewApp()
	app.Router.Run(app.Cfg.ServerAddress)
}