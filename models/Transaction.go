package models

import "time"

type Transaction struct {
	ID          int       `json:"id" db:"id"`
	FromAccount int       `json:"from_account" db:"from_account"`
	ToAccount   int       `json:"to_account" db:"to_account"`
	Amount      float64   `json:"amount" db:"amount"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Type        string    `json:"type" db:"type"`
}
