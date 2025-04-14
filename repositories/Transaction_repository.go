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

func (r *TransactionRepository) GetTransactionsByAccount(accountID int) ([]models.Transaction, error) {
	rows, err := r.db.Query(`
        SELECT id, from_account, to_account, amount, type, created_at 
        FROM transactions 
        WHERE from_account = $1 OR to_account = $1
        ORDER BY created_at DESC
    `, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(&t.ID, &t.FromAccount, &t.ToAccount, &t.Amount, &t.Type, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
