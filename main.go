package main

import (
	"bank_file_analyser/app/cmd"
	"bank_file_analyser/config"
	"bank_file_analyser/domain"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate)

	//load config
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	var app domain.App
	if conf.AppMode == "CMD" {
		app = cmd.NewCMDApp(conf)
	}
	if app == nil{
		log.Fatal("Couldn't start app. Check config.")
	}
	app.Run()
}
