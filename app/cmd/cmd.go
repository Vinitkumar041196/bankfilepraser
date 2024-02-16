package cmd

import (
	"bank_file_analyser/accounts/service"
	"bank_file_analyser/config"
	"bank_file_analyser/domain"
	"bank_file_analyser/fileparser"
	"fmt"
	"log"
	"os"
)

type cmdApp struct {
	Config            *config.AppConfig
	AccBalanceService domain.BalanceGeneratorService
}

func NewCMDApp(conf *config.AppConfig) domain.App {
	//initialize file parser
	parser := fileparser.NewCSVParser(rune(conf.FileColumnSeparator[0]), conf.FileHasHeader)

	//initialize accounts service
	accSrvc := service.NewBalanceGeneratorService(parser, conf.PayRefRegex, conf.DecimalPrecision)

	return &cmdApp{Config: conf, AccBalanceService: accSrvc}
}

func (app *cmdApp) Run() {
	file, err := os.Open(app.Config.FilePath)
	if err != nil {
		log.Fatal("Error while reading file: ", err)
	}
	defer file.Close()

	accBalances, err := app.AccBalanceService.GenerateAccBalancesFromFile(file)
	if err != nil {
		log.Fatal("Error while generating account balances: ", err)
	}

	resAccBalances := app.AccBalanceService.FormatAccountBalances(accBalances)
	fmt.Println("Totals")
	for curr, balance := range resAccBalances.Balances {
		fmt.Printf("%s %s\n", curr, balance.Total)
	}
}
