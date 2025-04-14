package models

import "time"

// BalancePrediction представляет прогноз баланса на определенную дату
type BalancePrediction struct {
	Date    time.Time `json:"date"`
	Balance float64   `json:"balance"`
	Type    string    `json:"type"` // starting_balance, payment, ending_balance
	Amount  float64   `json:"amount,omitempty"`
}
