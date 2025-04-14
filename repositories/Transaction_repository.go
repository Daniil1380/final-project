package repositories

import (
	"database/sql"
	"final-project/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	_, err := r.db.Exec(`
                INSERT INTO transactions (from_account, to_account, amount, type)
                VALUES ($1, $2, $3, $4)
        `, transaction.FromAccount, transaction.ToAccount, transaction.Amount, transaction.Type)
	return err
}
