package service

import (
	"bank_file_analyser/domain"
	"bank_file_analyser/utils"
	"io"
	"log"
	"time"
)

// Balances Generator Service
type balanceGeneratorService struct {
	paymentRefRegex  string
	decimalPrecision int
}

// initializes new balances generator service
func NewBalanceGeneratorService(payRefRegex string, decPrecision int) domain.BalanceGeneratorService {
	if payRefRegex == "" { //set default regex if input is empty
		payRefRegex = domain.PAYMENT_REFERENCE_REGEX
	}
	if decPrecision == 0 { //set default precision if input is empty
		decPrecision = domain.DECIMAL_PRECISION
	}
	return &balanceGeneratorService{paymentRefRegex: payRefRegex, decimalPrecision: decPrecision}
}

// Balances Generator
func (srvc *balanceGeneratorService) GenerateAccBalancesFromFile(parser domain.FileParser, file io.Reader, filterDateStr string) (*domain.BankAccBalances, error) {
	//parse file to struct
	accData, err := parser.Parse(file)
	if err != nil {
		log.Println("error while parsing file", err)
		return nil, err
	}

	var matchDate bool //by default dont match date

	//parse filter date 
	var filterDate time.Time
	if filterDateStr != "" {
		filterDate, err = time.Parse(domain.FILTER_DATE_FORMAT, filterDateStr)
		if err != nil {
			log.Println("error while parsing date", err)
			return nil, err
		}
		matchDate = true
	}

	//result variable
	accBalances := &domain.BankAccBalances{Balances: make(map[string]domain.AccBalances)}

	//loop through the file data and generate balances for each currency
	for _, row := range accData {
		//concat all naratives
		narrativeStr := row.Narrative1 + row.Narrative2 + row.Narrative3 + row.Narrative4 + row.Narrative5

		//match filter date
		if matchDate && !row.Date.Equal(filterDate) {
			continue
		}

		//check if narrative has payment reference
		if utils.MatchString(srvc.paymentRefRegex, narrativeStr) {
			amt := row.Debit - row.Credit //assuming amount can be in any column
			if amt < 0 {                  //using the absolute value
				amt *= -1
			}

			//add to the total for given currency
			bal := accBalances.Balances[row.Currency]
			bal.Total += amt
			accBalances.Balances[row.Currency] = bal
		}

	}

	return accBalances, nil
}

// used to format account balances for output
func (srvc *balanceGeneratorService) FormatAccountBalances(accBalances *domain.BankAccBalances) *domain.FormattedBankAccBalances {
	res := &domain.FormattedBankAccBalances{Balances: make(map[string]domain.FormattedAccBalances)}

	for curr, balance := range accBalances.Balances {
		res.Balances[curr] = domain.FormattedAccBalances{
			Total: utils.FormatInt64AmtToString(balance.Total, srvc.decimalPrecision),
		}
	}

	return res
}
