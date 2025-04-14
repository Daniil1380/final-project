package repositories

import (
	"database/sql"
	"final-project/models"
)

type CardRepository struct {
	DB *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{DB: db}
}

// Создание новой карты
func (r *CardRepository) CreateCard(card *models.Card) error {
	query := `INSERT INTO cards (account_id, encrypted_pan, cvv_hash, expiry_month, expiry_year, created_at) 
                  VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id`
	return r.DB.QueryRow(query, card.AccountID, card.EncryptedPAN, card.CVVHash, card.ExpiryMonth, card.ExpiryYear).Scan(&card.ID)
}

// Получение карты по ID
func (r *CardRepository) GetCardByID(id int) (*models.Card, error) {
	query := `SELECT id, account_id, encrypted_pan, cvv_hash, expiry_month, expiry_year, created_at FROM cards WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	var card models.Card
	err := row.Scan(&card.ID, &card.AccountID, &card.EncryptedPAN, &card.CVVHash, &card.ExpiryMonth, &card.ExpiryYear, &card.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &card, nil
}
