package app

import (
	"bank_file_analyser/accounts/service"
	"bank_file_analyser/config"
	_ "bank_file_analyser/docs"
	"bank_file_analyser/domain"
	"log"
	"strings"
)

// APP struct
type App struct {
	Config            *config.AppConfig
	AccBalanceService domain.BalanceGeneratorService
}

// create new app
func NewApp(conf *config.AppConfig) *App {
	//initialize accounts service
	accSrvc := service.NewBalanceGeneratorService(conf.PayRefRegex, conf.DecimalPrecision)

	return &App{Config: conf, AccBalanceService: accSrvc}
}

// start app
func (app *App) Start() {
	log.Println("Starting APP. Mode =", app.Config.AppMode)
	var err error
	//check App Mode
	switch strings.ToUpper(app.Config.AppMode) {
	case "CMD": // run using cli inputs
		err = RunCMDApp(app)
	case "HTTP": //start a http server
		err = RunHTTPApp(app)
	default:
		log.Fatal("Couldn't start app. Check config: APP_MODE")
	}

	if err != nil {
		log.Fatal("error while running app. ", err.Error())
	}
}
