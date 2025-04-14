package repositories

import (
	"database/sql"
	"final-project/models"
)

type PaymentScheduleRepository struct {
	db *sql.DB
}

func NewPaymentScheduleRepository(db *sql.DB) *PaymentScheduleRepository {
	return &PaymentScheduleRepository{db: db}
}

func (r *PaymentScheduleRepository) CreatePaymentSchedule(schedule *models.PaymentSchedule) error {
	_, err := r.db.Exec(`
        INSERT INTO payment_schedules (credit_id, amount, due_date, status)
        VALUES ($1, $2, $3, $4)
    `, schedule.CreditID, schedule.Amount, schedule.DueDate, schedule.Status)
	return err
}

func (r *PaymentScheduleRepository) GetPaymentSchedule(creditID int) ([]models.PaymentSchedule, error) {
	rows, err := r.db.Query(`
        SELECT id, credit_id, amount, due_date, status, paid_at 
        FROM payment_schedules 
        WHERE credit_id = $1
        ORDER BY due_date
    `, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PaymentSchedule
	for rows.Next() {
		var s models.PaymentSchedule
		err := rows.Scan(&s.ID, &s.CreditID, &s.Amount, &s.DueDate, &s.Status, &s.PaidAt)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}

	return schedules, nil
}
