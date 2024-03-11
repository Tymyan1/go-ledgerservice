package domain

import (
	"context"
	"database/sql"
)

// DbTransactionManager mock
type DbTransactionManagerMock struct {
	mockedBeginTxResult BeginTxResult
	beginTxMatcher      BeginTxMatcher
	beginTxCalls        []context.Context
}

type BeginTxResult struct {
	DbTransaction DbTransaction
	Error         error
}
type BeginTxMatcher struct {
	context context.Context
	options *sql.TxOptions
}

func dbTransactionManagerMock() DbTransactionManagerMock {
	return DbTransactionManagerMock{BeginTxResult{}, BeginTxMatcher{}, []context.Context{}}
}

func (mock *DbTransactionManagerMock) mockBeginTx(ctx context.Context, opts *sql.TxOptions, result DbTransaction, err error) {
	mock.mockedBeginTxResult = BeginTxResult{result, err}
	mock.beginTxMatcher = BeginTxMatcher{ctx, opts}
}

func (mock *DbTransactionManagerMock) BeginTx(ctx context.Context, opts *sql.TxOptions) (DbTransaction, error) {
	mock.beginTxCalls = append(mock.beginTxCalls, ctx)
	return mock.mockedBeginTxResult.DbTransaction, mock.mockedBeginTxResult.Error
}

// DbTransaction mock
type DbTransactionMock struct {
	mockedCommitResult   error
	commitCalls          int
	mockedRollbackResult error
	rollbackCalls        int
}

func dbTransactionMock() *DbTransactionMock {
	return &DbTransactionMock{}
}

func (mock *DbTransactionMock) mockCommit(result error) {
	mock.mockedCommitResult = result
}

func (mock *DbTransactionMock) mockRollback(result error) {
	mock.mockedRollbackResult = result
}

func (mock *DbTransactionMock) Commit() error {
	mock.commitCalls++
	return mock.mockedCommitResult
}

func (mock *DbTransactionMock) Rollback() error {
	mock.rollbackCalls++
	return mock.mockedRollbackResult
}
