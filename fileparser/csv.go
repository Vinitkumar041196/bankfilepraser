package fileparser

import (
	"bank_file_analyser/domain"
	"bank_file_analyser/utils"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

type CSVParser struct {
	separator        rune
	decimalPrecision int
	dateFormat       string
}

func NewCSVParser(separator rune, decimalPrecision int, dateFormat string) domain.FileParser {
	//setting default values if not available
	if separator == 0 {
		separator = ','
	}
	if decimalPrecision == 0 {
		decimalPrecision = domain.DECIMAL_PRECISION
	}
	if dateFormat == "" {
		dateFormat = domain.FILE_DATE_FORMAT
	}
	return &CSVParser{
		separator:        separator,
		decimalPrecision: decimalPrecision,
		dateFormat:       dateFormat,
	}
}

// buf can be file, byte buffer, etc
func (parser *CSVParser) Parse(buf io.Reader) ([]domain.BankStatementRecord, error) {
	if parser.separator == 0 {
		return nil, fmt.Errorf("couldn't get column separator")
	}

	r := csv.NewReader(buf)
	r.Comma = parser.separator

	//skip header
	r.Read()

	rows := []domain.BankStatementRecord{}
	for i := 0; ; i++ {
		line, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("error while parsing file", err)
			return nil, err
		}

		row := domain.BankStatementRecord{
			Narrative1: line[1],
			Narrative2: line[2],
			Narrative3: line[3],
			Narrative4: line[4],
			Narrative5: line[5],
			Type:       line[6],
			Currency:   strings.ToUpper(strings.TrimSpace(line[9])),
		}

		row.Date, err = time.Parse(parser.dateFormat, line[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing date in row %d: %w", i, err)
		}

		row.Credit, err = utils.FormatAmtStrToInt64(line[7], parser.decimalPrecision)
		if err != nil {
			log.Println("error while parsing number", err)
			return nil, fmt.Errorf("error parsing number in row %d: %w", i, err)
		}
		row.Debit, err = utils.FormatAmtStrToInt64(line[8], parser.decimalPrecision)
		if err != nil {
			log.Println("error while parsing number", err)
			return nil, fmt.Errorf("error parsing number in row %d: %w", i, err)
		}

		rows = append(rows, row)
	}

	return rows, nil
}
