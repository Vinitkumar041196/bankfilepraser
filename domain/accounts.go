package domain

import (
	"io"
	"time"
)

// BalanceGeneratorService
type BalanceGeneratorService interface {
	GenerateAccBalancesFromFile(parser FileParser, file io.Reader, filterDateStr string) (*BankAccBalances, error)
	FormatAccountBalances(accBalances *BankAccBalances) *FormattedBankAccBalances
}

// csv file data struct
type BankStatementRecord struct {
	Date       time.Time
	Narrative1 string
	Narrative2 string
	Narrative3 string
	Narrative4 string
	Narrative5 string
	Type       string
	Credit     int64
	Debit      int64
	Currency   string
}

// account balances per currency
type BankAccBalances struct {
	Balances map[string]AccBalances `json:"balances"`
}

// account balance fields
type AccBalances struct {
	Total int64 `json:"total"`
}

// used for output account balances per currency
type FormattedBankAccBalances struct {
	Balances map[string]FormattedAccBalances `json:"balances"`
}

// formatted values for balances
type FormattedAccBalances struct {
	Total string `json:"total"`
}
