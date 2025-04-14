package services

import (
	"errors"
	"final-project/models"
	"final-project/repositories"
	"final-project/utils"
	"math"
	"time"
)

type LoanService struct {
	LoanRepo    *repositories.LoanRepository
	AccountRepo *repositories.AccountRepository
}

func NewLoanService(loanRepo *repositories.LoanRepository, accountRepo *repositories.AccountRepository) *LoanService {
	return &LoanService{
		LoanRepo:    loanRepo,
		AccountRepo: accountRepo,
	}
}

// Расчет аннуитетного платежа
func (s *LoanService) CalculateMonthlyPayment(amount float64, term int, interestRate float64) float64 {
	// Преобразуем годовую ставку в месячную
	monthlyRate := interestRate / 12 / 100

	// Формула аннуитетного платежа
	payment := amount * (monthlyRate * math.Pow(1+monthlyRate, float64(term))) /
		(math.Pow(1+monthlyRate, float64(term)) - 1)

	return math.Round(payment*100) / 100 // Округляем до копеек
}

// Создание нового кредита
func (s *LoanService) CreateLoan(userID, accountID int, amount float64, term int) (*models.Loan, error) {
	// Проверяем, существует ли счет
	account, err := s.AccountRepo.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	// Проверяем, принадлежит ли счет пользователю
	if account.UserID != userID {
		return nil, errors.New("account does not belong to user")
	}

	interestRate, err := utils.GetCentralBankRate()

	// Рассчитываем ежемесячный платеж
	monthlyPayment := s.CalculateMonthlyPayment(amount, term, interestRate)

	loan := &models.Loan{
		UserID:          userID,
		AccountID:       accountID,
		Amount:          amount,
		Term:            term,
		InterestRate:    interestRate,
		MonthlyPayment:  monthlyPayment,
		Status:          "active",
		NextPaymentDate: time.Now().AddDate(0, 1, 0), // Первый платеж через месяц
	}

	// Сохраняем кредит в БД
	err = s.LoanRepo.CreateLoan(loan)
	if err != nil {
		return nil, err
	}

	// Зачисляем сумму кредита на счет
	err = s.AccountRepo.UpdateBalance(accountID, amount)
	if err != nil {
		utils.Logger.WithError(err).Error("Failed to deposit loan amount to account")
		return nil, err
	}

	return loan, nil
}

// ProcessCreditPayments проходит по активным кредитам, для которых наступила дата платежа,
// и списывает ежемесячный платёж (при наличии достаточного баланса на счёте).
// В случае недостатка средств начисляется штраф (+10% к сумме кредита).
func (s *LoanService) ProcessCreditPayments() error {
	// Предполагаем, что в репозитории LoanRepository реализован метод для получения просроченных платежей.
	loans, err := s.LoanRepo.GetActiveLoansDue(time.Now())
	if err != nil {
		return err
	}

	for _, loan := range loans {
		account, err := s.AccountRepo.GetAccountByID(loan.AccountID)
		if err != nil {
			continue
		}

		if account.Balance >= loan.MonthlyPayment {
			// Списываем платёж
			err = s.AccountRepo.UpdateBalance(loan.AccountID, -loan.MonthlyPayment)
			if err != nil {
				continue
			}
			// Обновляем дату следующего платежа (прибавляем 1 месяц)
			loan.NextPaymentDate = loan.NextPaymentDate.AddDate(0, 1, 0)
			// Если кредит погашается, можно обновить статус, здесь оставляем статус "active"
			err = s.LoanRepo.UpdateLoan(loan)
			if err != nil {
				continue
			}
		} else {
			// Недостаточно средств — начисляем штраф 10%
			penalty := loan.MonthlyPayment * 0.10
			loan.Amount = math.Round((loan.Amount+penalty)*100) / 100
			err = s.LoanRepo.UpdateLoan(loan)
			if err != nil {
				continue
			}
		}
	}
	return nil
}
