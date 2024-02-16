package domain

import "io"

type FileParser interface {
	Parse(buf io.Reader) ([]BankStatementRecord, error)
}
