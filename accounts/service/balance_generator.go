package service

import (
	"bank_file_analyser/domain"
	"bank_file_analyser/utils"
	"io"
	"log"
)

type balanceGeneratorService struct {
	paymentRefRegex    string
	decimalPrecision int
	fileParser       domain.FileParser
}

func NewBalanceGeneratorService(parser domain.FileParser, payRefRegex string, decPrecision int) domain.BalanceGeneratorService {
	if payRefRegex == "" {
		payRefRegex = domain.PAYMENT_REFERENCE_REGEX
	}
	if decPrecision == 0 {
		decPrecision = domain.DECIMAL_PRECISION
	}
	return &balanceGeneratorService{fileParser: parser, paymentRefRegex: payRefRegex, decimalPrecision: decPrecision}
}

func (srvc *balanceGeneratorService) FormatAccountBalances(accBalances *domain.BankAccBalances) (*domain.FormattedBankAccBalances) {
	res := &domain.FormattedBankAccBalances{Balances: make(map[string]domain.FormattedAccBalances)}

	for curr, balance := range accBalances.Balances {
		res.Balances[curr] = domain.FormattedAccBalances{
			Total: utils.FormatInt64AmtToString(balance.Total, srvc.decimalPrecision),
		}
	}

	return res
}

func (srvc *balanceGeneratorService) GenerateAccBalancesFromFile(file io.Reader) (*domain.BankAccBalances, error) {
	//parse file to struct
	accData, err := srvc.fileParser.Parse(file)
	if err != nil {
		log.Println("Error while parsing file", err)
		return nil, err
	}

	//result
	accBalances := &domain.BankAccBalances{Balances: make(map[string]domain.AccBalances)}

	for _, row := range accData {
		//concat all naratives
		narrativeStr := row.Narrative1 + row.Narrative2 + row.Narrative3 + row.Narrative4 + row.Narrative5

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
