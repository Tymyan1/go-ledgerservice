package domain

import (
	"context"
	"database/sql"
)

// leaking sql-specific tx options for convenience over declaring yet another interface, especially since this is currently always nil
type DbTransactionManager interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (DbTransaction, error)
}

type DbTransaction interface {
	Rollback() error
	Commit() error
}
