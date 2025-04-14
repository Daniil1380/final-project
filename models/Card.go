package models

import "time"

// Card представляет виртуальную карту пользователя
type Card struct {
	ID           int       `json:"id"`
	AccountID    int       `json:"account_id"`
	EncryptedPAN string    `json:"-"` // Зашифрованный номер карты
	CVVHash      string    `json:"-"` // Захешированный CVV
	ExpiryMonth  int       `json:"expiry_month"`
	ExpiryYear   int       `json:"expiry_year"`
	CreatedAt    time.Time `json:"created_at"`
}
