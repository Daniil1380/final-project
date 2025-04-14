package models

import "time"

type PaymentSchedule struct {
	ID       int        `json:"id" db:"id"`
	CreditID int        `json:"credit_id" db:"credit_id"`
	Amount   float64    `json:"amount" db:"amount"`
	DueDate  time.Time  `json:"due_date" db:"due_date"`
	Status   string     `json:"status" db:"status"` // "pending", "paid", "overdue"
	PaidAt   *time.Time `json:"paid_at,omitempty" db:"paid_at"`
}
