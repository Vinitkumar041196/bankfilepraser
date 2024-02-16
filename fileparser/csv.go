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
	separator rune
	hasHeader bool
}

func NewCSVParser(separator rune, hasHeader bool) domain.FileParser {
	if separator == 0 {
		separator = ','
	}
	return &CSVParser{
		separator: separator,
		hasHeader: hasHeader,
	}
}

// buf can be file, byte buffer, etc
func (tp *CSVParser) Parse(buf io.Reader) ([]domain.BankStatementRecord, error) {
	if tp.separator == 0 {
		return nil, fmt.Errorf("couldn't get column separator")
	}

	r := csv.NewReader(buf)
	r.Comma = tp.separator

	if tp.hasHeader { //skip header
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
		row.Credit, err = utils.FormatAmtStrToInt64(line[7], domain.DECIMAL_PRECISION)
		if err != nil {
			log.Println("Error while parsing numbers", err)
			return nil, err
		}
		row.Debit, err = utils.FormatAmtStrToInt64(line[8], domain.DECIMAL_PRECISION)
		if err != nil {
			log.Println("Error while parsing numbers", err)
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, nil
}
