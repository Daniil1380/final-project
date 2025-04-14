package services

import (
	"errors"
	"final-project/models"
	"final-project/repositories"
)

type AccountService struct {
	AccountRepo *repositories.AccountRepository
}

func NewAccountService(accountRepo *repositories.AccountRepository) *AccountService {
	return &AccountService{AccountRepo: accountRepo}
}

// Создание нового счета
func (s *AccountService) CreateAccount(userID int, currency string) (*models.Account, error) {
	account := &models.Account{
		UserID:   userID,
		Balance:  0.0,
		Currency: currency,
	}

	err := s.AccountRepo.CreateAccount(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// Пополнение счета
func (s *AccountService) Deposit(accountID int, amount float64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}
	return s.AccountRepo.UpdateBalance(accountID, amount)
}

// Списание средств
func (s *AccountService) Withdraw(accountID int, amount float64) error {
	account, err := s.AccountRepo.GetAccountByID(accountID)
	if err != nil {
		return err
	}

	if account.Balance < amount {
		return errors.New("insufficient funds")
	}

	return s.AccountRepo.UpdateBalance(accountID, -amount)
}

func (s *AccountService) TransferFunds(fromAccountID, toAccountID int, amount float64, userID int) error {
	// Проверяем, что сумма положительная
	if amount <= 0 {
		return errors.New("transfer amount must be positive")
	}

	// Проверяем, принадлежит ли счет пользователю
	fromAccount, err := s.AccountRepo.GetAccountByID(fromAccountID)
	if err != nil {
		return err
	}

	if fromAccount.UserID != userID {
		return errors.New("account does not belong to user")
	}

	// Проверяем существование счета получателя
	_, err = s.AccountRepo.GetAccountByID(toAccountID)
	if err != nil {
		return errors.New("recipient account not found")
	}

	// Выполняем перевод через репозиторий
	return s.AccountRepo.TransferFunds(fromAccountID, toAccountID, amount)
}
