package domain

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

const (
	accountKey = "accountKey1"
	timestamp1 = "2021-01-01"
	timestamp2 = "2021-01-02"
)

func TestTransactionService_CreateTransaction_GoldenPath(t *testing.T) {
	transactionService, transactionPortMock, dbTransactionManagerMock := setupService()
	dbTransaction := dbTransactionMock()
	latestTransaction := Transaction{PostedTimestamp: parse(timestamp1), Balance: decimal.NewFromInt(10)}
	ctx := context.TODO()
	transaction := Transaction{Amount: decimal.NewFromInt(10), AccountKey: accountKey, Type: CREDIT, PostedTimestamp: parse(timestamp2)}
	transactionWithBalance := transaction
	transactionWithBalance.Balance = decimal.NewFromInt(20)

	dbTransactionManagerMock.mockBeginTx(ctx, nil, dbTransaction, nil)
	transactionPortMock.mockQueryLatestForAccount(accountKey, &latestTransaction, nil)
	transactionPortMock.mockSave(transactionWithBalance, nil)
	dbTransaction.mockCommit(nil)

	err := transactionService.ProcessTransaction(transaction, ctx)

	assert.Nil(t, err, "there should not be an error")
	assert.Equal(t, 1, len(dbTransactionManagerMock.beginTxCalls), "beginTx should be called")
	assert.Equal(t, accountKey, transactionPortMock.queryLatestForAccountCalls[0], "queryLatestForAccount should be called with the account key")
	assert.Equal(t, transactionWithBalance, transactionPortMock.saveCalls[0], "save should be called with the transaction and correct balance")
	assert.Equal(t, 1, dbTransaction.commitCalls, "commit should be called")
}

func TestTransactionService_CreateTransaction_ShouldFail_ForDebitTransaction_WhenNoPreviousTransactionExists(t *testing.T) {
	transactionService, transactionPortMock, dbTransactionManagerMock := setupService()
	dbTransaction := dbTransactionMock()

	ctx := context.TODO()
	transaction := Transaction{Amount: decimal.NewFromInt(10), AccountKey: accountKey, Type: DEBIT, PostedTimestamp: parse(timestamp2)}

	dbTransactionManagerMock.mockBeginTx(ctx, nil, dbTransaction, nil)
	transactionPortMock.mockQueryLatestForAccount(accountKey, nil, nil)
	dbTransaction.mockRollback(nil)

	err := transactionService.ProcessTransaction(transaction, ctx)

	assert.Equal(t, "insufficient funds", err.Error(), "an error is returned")
	assert.Equal(t, 1, len(dbTransactionManagerMock.beginTxCalls), "beginTx should be called")
	assert.Equal(t, accountKey, transactionPortMock.queryLatestForAccountCalls[0], "queryLatestForAccount should be called with the account key")
	assert.Equal(t, 0, len(transactionPortMock.saveCalls), "save should not be called")
	assert.Equal(t, 0, dbTransaction.commitCalls, "commit should not be called")
	assert.Equal(t, 1, dbTransaction.rollbackCalls, "rollback should be called")
}

func TestTransactionService_CreateTransaction_ShouldFail_WhenDebitTransaction_AndInsufficientBalance(t *testing.T) {
	transactionService, transactionPortMock, dbTransactionManagerMock := setupService()
	dbTransaction := dbTransactionMock()
	latestTransaction := Transaction{PostedTimestamp: parse(timestamp1), Balance: decimal.NewFromInt(10)}
	ctx := context.TODO()
	transaction := Transaction{Amount: decimal.NewFromInt(100), AccountKey: accountKey, Type: DEBIT, PostedTimestamp: parse(timestamp2)}

	dbTransactionManagerMock.mockBeginTx(ctx, nil, dbTransaction, nil)
	transactionPortMock.mockQueryLatestForAccount(accountKey, &latestTransaction, nil)
	dbTransaction.mockRollback(nil)

	err := transactionService.ProcessTransaction(transaction, ctx)

	assert.Equal(t, "insufficient funds", err.Error(), "an error is returned")
	assert.Equal(t, 1, len(dbTransactionManagerMock.beginTxCalls), "beginTx should be called")
	assert.Equal(t, accountKey, transactionPortMock.queryLatestForAccountCalls[0], "queryLatestForAccount should be called with the account key")
	assert.Equal(t, 0, len(transactionPortMock.saveCalls), "save should not be called")
	assert.Equal(t, 1, dbTransaction.rollbackCalls, "rollback should be called")
}

func setupService() (*TransactionService, *TransactionPortMock, *DbTransactionManagerMock) {
	transactionPortMock := transactionPortMock()
	dbTransactionManagerMock := dbTransactionManagerMock()
	transactionService := TransactionService{&transactionPortMock, &dbTransactionManagerMock}
	return &transactionService, &transactionPortMock, &dbTransactionManagerMock
}

func transactionPortMock() TransactionPortMock {
	return TransactionPortMock{QueryLatestForAccountResult{}, "", []string{}, nil, nil, []Transaction{}}
}

// TransactionPort mock
type TransactionPortMock struct {
	mockedQueryLatestForAccountResult QueryLatestForAccountResult
	queryLatestForAccountMatcher      string
	queryLatestForAccountCalls        []string
	mockedSaveResult                  error
	saveMatcher                       *Transaction
	saveCalls                         []Transaction
}

type QueryLatestForAccountResult struct {
	Transaction *Transaction
	Error       error
}

func (mock *TransactionPortMock) Save(transaction Transaction) error {
	mock.saveCalls = append(mock.saveCalls, transaction)
	if reflect.DeepEqual(transaction, *mock.saveMatcher) {
		return mock.mockedSaveResult
	}
	panic(fmt.Sprintf("save called with unmatched transaction %s, expected transaction %s", transaction, *mock.saveMatcher))
}

func (tmp *TransactionPortMock) QueryLatestForAccount(accountKey string) (*Transaction, error) {
	tmp.queryLatestForAccountCalls = append(tmp.queryLatestForAccountCalls, accountKey)
	if accountKey == tmp.queryLatestForAccountMatcher {
		return tmp.mockedQueryLatestForAccountResult.Transaction, tmp.mockedQueryLatestForAccountResult.Error
	}
	panic(fmt.Sprintf("QueryLatestForAccount called with unmatched accountKey %s, expected accountKey %s", accountKey, tmp.queryLatestForAccountMatcher))
}

func (mock *TransactionPortMock) mockSave(matcher Transaction, result error) {
	mock.saveMatcher = &matcher
	mock.mockedSaveResult = result
}

func (mock *TransactionPortMock) mockQueryLatestForAccount(matcher string, result *Transaction, err error) {
	mock.queryLatestForAccountMatcher = matcher
	mock.mockedQueryLatestForAccountResult = QueryLatestForAccountResult{Transaction: result, Error: err}
}

func parse(s string) time.Time {
	t, _ := time.Parse(time.DateOnly, s)
	return t
}
