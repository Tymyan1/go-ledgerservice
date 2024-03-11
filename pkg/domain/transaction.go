package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	Key             string
	Amount          decimal.Decimal
	Type            TransactionType
	AccountKey      string
	Balance         decimal.Decimal
	PostedTimestamp time.Time
}

type TransactionType string

const (
	CREDIT TransactionType = "CREDIT"
	DEBIT  TransactionType = "DEBIT"
)
