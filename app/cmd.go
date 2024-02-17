package app

import (
	"bank_file_analyser/fileparser"
	"flag"
	"fmt"
	"log"
	"os"
)

// CMD version of app
func RunCMDApp(app *App) error {
	//initialize file parser
	parser := fileparser.NewCSVParser(rune(app.Config.FileColumnSeparator[0]), app.Config.DecimalPrecision)

	filePath := flag.String("file_path", "", "path to statement file")
	filterDate := flag.String("date", "", "filter date format: DD/MM/YYYY")
	outFilePath := flag.String("out_file_path", "", "if provided the result will be stored in file.")
	flag.Parse()

	if *filePath == "" {
		return fmt.Errorf("no file provided to process")
	}

	file, err := os.Open(*filePath)
	if err != nil {
		log.Println("error while reading file: ", err)
		return nil
	}
	defer file.Close()

	accBalances, err := app.AccBalanceService.GenerateAccBalancesFromFile(parser, file, *filterDate)
	if err != nil {
		log.Println("error while generating account balances: ", err)
		return nil
	}

	resAccBalances := app.AccBalanceService.FormatAccountBalances(accBalances)

	if *outFilePath == "" {
		fmt.Println("\nResult:\nCurrency | Totals")
		for curr, balance := range resAccBalances.Balances {
			fmt.Printf("%8s | %s\n", curr, balance.Total)
		}
	} else {
		resFile, err := os.Create(*outFilePath)
		if err != nil {
			log.Println("error opening output file:", err)
			return err
		}
		defer resFile.Close()
		fmt.Fprintf(resFile, "Currency,Totals\n")
		for curr, balance := range resAccBalances.Balances {
			fmt.Fprintf(resFile, "%s,%s\n", curr, balance.Total)
		}
	}

	return nil
}
