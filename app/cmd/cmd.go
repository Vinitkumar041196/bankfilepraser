package cmd

import (
	"bank_file_analyser/accounts/service"
	"bank_file_analyser/config"
	"bank_file_analyser/domain"
	"bank_file_analyser/fileparser"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
	c := make(chan os.Signal, 1)
	defer close(c)

	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Stopping.")
		os.Exit(1)
	}()

	fmt.Println("\nStarting Statement Processor.\nPress ctrl+c to exit.")
	for {
		//get file path
		fmt.Println("\nEnter file path:")
		filePath := ""
		fmt.Scanln(&filePath)

		filePath = strings.TrimSpace(filePath)
		if filePath == "" {
			continue
		}

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error while reading file: ", err)
			continue
		}

		accBalances, err := app.AccBalanceService.GenerateAccBalancesFromFile(file)
		if err != nil {
			fmt.Println("Error while generating account balances: ", err)
			continue
		}

		file.Close()

		resAccBalances := app.AccBalanceService.FormatAccountBalances(accBalances)

		fmt.Println("\nResult:\nCurrency | Totals")
		for curr, balance := range resAccBalances.Balances {
			fmt.Printf("%8s | %s\n", curr, balance.Total)
		}

		fmt.Println("\nDo you have more files to process? Y/N: ")

		more := "N"
		fmt.Scanln(&more)
		if strings.ToUpper(more) == "N" {
			fmt.Println("Stopping")
			break
		}
	}
}
