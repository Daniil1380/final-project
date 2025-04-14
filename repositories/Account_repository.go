package repositories

import (
	"database/sql"
	"final-project/models"
)

type AccountRepository struct {
	DB *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{DB: db}
}

// Создание нового счета
func (r *AccountRepository) CreateAccount(account *models.Account) error {
	query := `INSERT INTO accounts (user_id, balance, currency, created_at) 
                  VALUES ($1, $2, $3, NOW()) RETURNING id`
	return r.DB.QueryRow(query, account.UserID, account.Balance, account.Currency).Scan(&account.ID)
}

// Получение счета по ID
func (r *AccountRepository) GetAccountByID(id int) (*models.Account, error) {
	query := `SELECT id, user_id, balance, currency, created_at FROM accounts WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	var account models.Account
	err := row.Scan(&account.ID, &account.UserID, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// Обновление баланса счета
func (r *AccountRepository) UpdateBalance(accountID int, amount float64) error {
	query := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	_, err := r.DB.Exec(query, amount, accountID)
	return err
}
