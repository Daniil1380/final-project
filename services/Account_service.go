package services

import (
	"errors"
	"final-project/models"
	"final-project/repositories"
	"final-project/utils"
	"log"
	"strconv"
	"time"
)

type AccountService struct {
	AccountRepo     *repositories.AccountRepository
	TransactionRepo *repositories.TransactionRepository
}

func NewAccountService(accountRepo *repositories.AccountRepository, transactionRepo *repositories.TransactionRepository) *AccountService {
	return &AccountService{
		AccountRepo:     accountRepo,
		TransactionRepo: transactionRepo,
	}
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

	// Обновление баланса
	err := s.AccountRepo.UpdateBalance(accountID, amount)
	if err != nil {
		return err
	}

	// Сохранение транзакции
	transaction := &models.Transaction{
		ID:        accountID,
		Amount:    amount,
		Type:      "deposit",
		CreatedAt: time.Now(),
	}

	return s.TransactionRepo.CreateTransaction(transaction)
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

	// Обновление баланса
	err = s.AccountRepo.UpdateBalance(accountID, -amount)
	if err != nil {
		return err
	}

	// Сохранение транзакции
	transaction := &models.Transaction{
		ID:        accountID,
		Amount:    -amount,
		Type:      "withdraw",
		CreatedAt: time.Now(),
	}

	return s.TransactionRepo.CreateTransaction(transaction)
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

	var account *models.Account

	// Проверяем существование счета получателя
	account, err = s.AccountRepo.GetAccountByID(toAccountID)
	if err != nil {
		return errors.New("recipient account not found")
	}

	// Вызываем нашу функцию для отправки email
	err = utils.SendPaymentEmail(strconv.Itoa(account.UserID), amount)
	if err != nil {
		// Если функция вернула ошибку, выводим фатальное сообщение и завершаем программу
		// Логирование самой ошибки уже произошло внутри sendEmail и sendPaymentEmail
		log.Fatalf("Не удалось отправить email: %v", err)
	}

	log.Println("Процесс отправки тестового email завершен успешно.")

	// Выполняем перевод через репозиторий
	err = s.AccountRepo.TransferFunds(fromAccountID, toAccountID, amount)
	if err != nil {
		return err
	}

	// Сохранение транзакций
	// Транзакция списания
	outTransaction := &models.Transaction{
		ID:        fromAccountID,
		Amount:    -amount,
		Type:      "transfer_out",
		CreatedAt: time.Now(),
	}

	err = s.TransactionRepo.CreateTransaction(outTransaction)
	if err != nil {
		return err
	}

	// Транзакция зачисления
	inTransaction := &models.Transaction{
		ID:        toAccountID,
		Amount:    amount,
		Type:      "transfer_in",
		CreatedAt: time.Now(),
	}

	return s.TransactionRepo.CreateTransaction(inTransaction)
}
