package fileparser

import (
	"bank_file_analyser/domain"
	"bank_file_analyser/utils"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
)

type CSVParser struct {
	separator        rune
	hasHeader        bool
	decimalPrecision int
}

func NewCSVParser(separator rune, hasHeader bool, decimalPrecision int) domain.FileParser {
	if separator == 0 {
		separator = ','
	}
	if decimalPrecision == 0 {
		decimalPrecision = domain.DECIMAL_PRECISION
	}
	return &CSVParser{
		separator:        separator,
		hasHeader:        hasHeader,
		decimalPrecision: decimalPrecision,
	}
}

// buf can be file, byte buffer, etc
func (parser *CSVParser) Parse(buf io.Reader) ([]domain.BankStatementRecord, error) {
	if parser.separator == 0 {
		return nil, fmt.Errorf("couldn't get column separator")
	}

	r := csv.NewReader(buf)
	r.Comma = parser.separator

	if parser.hasHeader { //skip header
		r.Read()
	}

	rows := []domain.BankStatementRecord{}
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error while parsing file", err)
			return nil, err
		}

		row := domain.BankStatementRecord{
			Date:       line[0],
			Narrative1: line[1],
			Narrative2: line[2],
			Narrative3: line[3],
			Narrative4: line[4],
			Narrative5: line[5],
			Type:       line[6],
			Currency:   strings.ToUpper(strings.TrimSpace(line[9])),
		}
		row.Credit, err = utils.FormatAmtStrToInt64(line[7], parser.decimalPrecision)
		if err != nil {
			log.Println("Error while parsing numbers", err)
			return nil, err
		}
		row.Debit, err = utils.FormatAmtStrToInt64(line[8], parser.decimalPrecision)
		if err != nil {
			log.Println("Error while parsing numbers", err)
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, nil
}
