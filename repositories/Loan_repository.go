// LoanRepository.go

package repositories

import (
	"database/sql"
	"final-project/models"
	"final-project/utils"
)

// LoanRepository представляет собой репозиторий для работы с кредитами
type LoanRepository struct {
	db *sql.DB
}

// NewLoanRepository создает новый экземпляр LoanRepository
func NewLoanRepository(db *sql.DB) *LoanRepository {
	return &LoanRepository{
		db: db,
	}
}

// CreateLoan создает новый кредит в базе данных
func (r *LoanRepository) CreateLoan(loan *models.Loan) error {
	query := `
                INSERT INTO loans (
                        user_id,
                        account_id,
                        amount,
                        term,
                        interest_rate,
                        monthly_payment,
                        status,
                        created_at,
                        next_payment_date
                ) VALUES (
                        $1,
                        $2,
                        $3,
                        $4,
                        $5,
                        $6,
                        $7,
                        $8,
                        $9
                ) RETURNING id
        `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to prepare query for creating loan")
		return err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(
		loan.UserID,
		loan.AccountID,
		loan.Amount,
		loan.Term,
		loan.InterestRate,
		loan.MonthlyPayment,
		loan.Status,
		loan.CreatedAt,
		loan.NextPaymentDate,
	).Scan(&id)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to create loan")
		return err
	}

	loan.ID = id

	return nil
}

// GetLoanByID получает кредит по его идентификатору
func (r *LoanRepository) GetLoanByID(id int) (*models.Loan, error) {
	query := `
                SELECT 
                        id,
                        user_id,
                        account_id,
                        amount,
                        term,
                        interest_rate,
                        monthly_payment,
                        status,
                        created_at,
                        next_payment_date
                FROM loans
                WHERE id = $1
        `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to prepare query for getting loan by ID")
		return nil, err
	}
	defer stmt.Close()

	var loan models.Loan
	err = stmt.QueryRow(id).Scan(
		&loan.ID,
		&loan.UserID,
		&loan.AccountID,
		&loan.Amount,
		&loan.Term,
		&loan.InterestRate,
		&loan.MonthlyPayment,
		&loan.Status,
		&loan.CreatedAt,
		&loan.NextPaymentDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		utils.Logger.WithError(err).Error("Failed to get loan by ID")
		return nil, err
	}

	return &loan, nil
}

// UpdateLoan обновляет существующий кредит
func (r *LoanRepository) UpdateLoan(loan *models.Loan) error {
	query := `
                UPDATE loans
                SET 
                        user_id = $1,
                        account_id = $2,
                        amount = $3,
                        term = $4,
                        interest_rate = $5,
                        monthly_payment = $6,
                        status = $7,
                        created_at = $8,
                        next_payment_date = $9
                WHERE id = $10
        `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to prepare query for updating loan")
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		loan.UserID,
		loan.AccountID,
		loan.Amount,
		loan.Term,
		loan.InterestRate,
		loan.MonthlyPayment,
		loan.Status,
		loan.CreatedAt,
		loan.NextPaymentDate,
		loan.ID,
	)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to update loan")
		return err
	}

	return nil
}

// DeleteLoan удаляет кредит по его идентификатору
func (r *LoanRepository) DeleteLoan(id int) error {
	query := `
                DELETE FROM loans
                WHERE id = $1
        `

	stmt, err := r.db.Prepare(query)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to prepare query for deleting loan")
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to delete loan")
		return err
	}

	return nil
}
