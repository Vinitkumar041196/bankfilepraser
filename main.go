package main

import (
	"bank_file_analyser/app/cmd"
	"bank_file_analyser/app/http"
	"bank_file_analyser/config"
	"bank_file_analyser/domain"
	"log"
)

// @title Statement Processor
// @version 1.0
// @BasePath /v1

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate)

	//load config
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	var app domain.App
	switch conf.AppMode {
	case "CMD":
		app = cmd.NewCMDApp(conf)
	case "HTTP":
		app = http.NewHttpApp(conf)
	default:
		log.Fatal("Couldn't start app. Check config.")
	}

	app.Run()
}
