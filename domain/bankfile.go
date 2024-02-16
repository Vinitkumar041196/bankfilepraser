package domain

const (
	PAYMENT_REFERENCE_TAG = "PAY[0-9]{6}[a-zA-Z]{2}"
)

type BankStatementRecord struct {
	Date       string
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

type BankAccBalances struct {
	Balances map[string]AccBalances
}

type AccBalances struct {
	Total int64
}
