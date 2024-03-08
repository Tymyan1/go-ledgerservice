package domain

import (
	"github.com/shopspring/decimal"
)

type Transaction struct {
	Key        string
	Amount     decimal.Decimal
	Type       TransactionType
	AccountKey string
}

type TransactionType string

const (
	CREDIT TransactionType = "CREDIT"
	DEBIT  TransactionType = "DEBIT"
)
