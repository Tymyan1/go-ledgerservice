package api

import (
	"errors"
	"ledgerservice/pkg/domain"
	"time"

	"github.com/shopspring/decimal"
)

type TransactionDto struct {
	Key             *string
	Amount          *decimal.Decimal
	Type            *domain.TransactionType
	AccountKey      *string
	PostedTimestamp *time.Time
}

func validate(dto TransactionDto) []error {
	var errList []error
	if dto.Key == nil {
		errList = append(errList, errors.New("field key is required"))
	}

	if dto.Amount == nil {
		errList = append(errList, errors.New("field amount is required"))
	} else if dto.Amount.LessThan(decimal.Zero) {
		errList = append(errList, errors.New("field amount must not be negative"))
	}

	if dto.Type == nil {
		errList = append(errList, errors.New("field type is required"))
	} else if *dto.Type != domain.CREDIT && *dto.Type != domain.DEBIT {
		errList = append(errList, errors.New("field type must be CREDIT or DEBIT"))
	}

	if dto.AccountKey == nil {
		errList = append(errList, errors.New("field accountKey is required"))
	}

	if dto.PostedTimestamp == nil {
		errList = append(errList, errors.New("field postedTimestamp is required"))
	}
	return errList
}

func mapTransactionToDomain(transaction TransactionDto) domain.Transaction {
	return domain.Transaction{
		Key:             *transaction.Key,
		Amount:          *transaction.Amount,
		Type:            *transaction.Type,
		AccountKey:      *transaction.AccountKey,
		PostedTimestamp: *transaction.PostedTimestamp,
	}
}
