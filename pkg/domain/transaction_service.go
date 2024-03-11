package domain

import (
	"context"
	"errors"
	"log"

	"github.com/shopspring/decimal"
)

type TransactionService struct {
	transactionDb TransactionPort
	txManager     DbTransactionManager
}

type TransactionPort interface {
	Save(transaction Transaction) error
	QueryLatestForAccount(accountKey string) (*Transaction, error)
}

func (ts *TransactionService) ProcessTransaction(transaction Transaction, ctx context.Context) error {
	dbTx, err := ts.txManager.BeginTx(ctx, nil)
	if err != nil {
		return rollback(dbTx, err)
	}

	lastTxPt, err := ts.transactionDb.QueryLatestForAccount(transaction.AccountKey)
	if err != nil {
		return rollback(dbTx, err)
	}

	lastTx := getTransactionData(lastTxPt)

	if err := validateTransaction(transaction, lastTx); err != nil {
		return rollback(dbTx, err)
	}

	newBalance, err := calculateNewBalance(transaction, lastTx)
	if err != nil {
		return rollback(dbTx, err)
	}

	transaction.Balance = newBalance
	if err := ts.transactionDb.Save(transaction); err != nil {
		return rollback(dbTx, err)
	}

	return dbTx.Commit()
}

func rollback(tx DbTransaction, err error) error {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}
	return err
}

func calculateNewBalance(newTransaction Transaction, lastTransaction Transaction) (decimal.Decimal, error) {
	switch newTransaction.Type {
	case CREDIT:
		return lastTransaction.Balance.Add(newTransaction.Amount), nil
	case DEBIT:
		return lastTransaction.Balance.Sub(newTransaction.Amount), nil
	}
	return decimal.Zero, errors.New("invalid transaction type")
}

func validateTransaction(newTransaction Transaction, lastTransaction Transaction) error {
	if lastTransaction.PostedTimestamp.After(newTransaction.PostedTimestamp) {
		return errors.New("latest processed transaction is posted after this transaction")
	}
	if lastTransaction.PostedTimestamp.Equal(newTransaction.PostedTimestamp) {
		return errors.New("the posted timestamp must be unique")
	}
	if newTransaction.Type == DEBIT && newTransaction.Amount.GreaterThan(lastTransaction.Balance) {
		return errors.New("insufficient funds")
	}

	return nil
}

func getTransactionData(txPt *Transaction) Transaction {
	if txPt == nil {
		log.Println("Using implicit 0 balance")
		return implicitZeroBalanceTransaction()
	}
	return *txPt
}

func implicitZeroBalanceTransaction() Transaction {
	return Transaction{
		Balance: decimal.Zero,
	}
}
