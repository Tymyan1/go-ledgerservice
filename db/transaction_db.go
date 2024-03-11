package db

import (
	"database/sql"
	"ledgerservice/pkg/domain"
)

type TransactionDb struct {
	DB *sql.DB
}

func (db *TransactionDb) Save(transaction domain.Transaction) error {
	_, err := db.DB.Exec("INSERT INTO transactions (key, amount, type, account_key, balance, posted_timestamp) VALUES ($1, $2, $3, $4, $5, $6)",
		transaction.Key, transaction.Amount, transaction.Type, transaction.AccountKey, transaction.Balance, transaction.PostedTimestamp)
	return err
}

func (db *TransactionDb) QueryAll() ([]domain.Transaction, error) {
	rows, err := db.DB.Query("SELECT * FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []domain.Transaction

	for rows.Next() {
		var tx domain.Transaction

		err := rows.Scan(&tx.Key, &tx.Amount, &tx.Type, &tx.AccountKey, &tx.Balance, &tx.PostedTimestamp)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, tx)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
