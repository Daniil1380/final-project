package services

import (
	"final-project/models"
	"final-project/repositories"
)

type TransactionService struct {
	TransactionRepo *repositories.TransactionRepository
}

func NewTransactionService(transactionRepo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{TransactionRepo: transactionRepo}
}

func (s *TransactionService) CreateTransaction(fromAccountID, toAccountID int, amount float64, transactionType string) error {
	transaction := &models.Transaction{
		FromAccount: fromAccountID,
		ToAccount:   toAccountID,
		Amount:      amount,
		Type:        transactionType,
	}

	return s.TransactionRepo.CreateTransaction(transaction)
}
