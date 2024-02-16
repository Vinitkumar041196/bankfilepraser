package main

import (
	"bank_file_analyser/accounts/service"
	"bank_file_analyser/config"
	"bank_file_analyser/fileparser"
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate)

	//load config
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	//initialize file parser
	parser := fileparser.NewCSVParser(rune(conf.FileColumnSeparator[0]), conf.FileHasHeader)

	//initialize accounts service
	accSrvc := service.NewBalanceGeneratorService(parser, conf.PayRefRegex, conf.DecimalPrecision)

	file, err := os.Open(conf.FilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := accSrvc.GenerateAccBalancesFromFile(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(accSrvc.FormatAccountBalances(data))
}
