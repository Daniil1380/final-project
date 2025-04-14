package models

import "time"

type Loan struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	AccountID       int       `json:"account_id"`
	Amount          float64   `json:"amount"`
	Term            int       `json:"term"`          // в месяцах
	InterestRate    float64   `json:"interest_rate"` // годовая процентная ставка
	MonthlyPayment  float64   `json:"monthly_payment"`
	Status          string    `json:"status"` // active, paid, defaulted
	CreatedAt       time.Time `json:"created_at"`
	NextPaymentDate time.Time `json:"next_payment_date"`
}
