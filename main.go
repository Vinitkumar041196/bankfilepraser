package main

import (
	"bank_file_analyser/app"
	"bank_file_analyser/config"
	"log"
)

// @title Statement Processor
// @version 1.0
// @BasePath /v1
func main() {
	//set log flags
	log.SetFlags(log.Lshortfile | log.Ldate)

	//load config
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	//new app
	app := app.NewApp(conf)
	//start app
	app.Start()
}
