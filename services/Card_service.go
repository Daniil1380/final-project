package services

import (
	"crypto/sha256"
	"encoding/hex"
	_ "errors"
	"final-project/models"
	"final-project/repositories"
	"golang.org/x/crypto/bcrypt"
)

type CardService struct {
	CardRepo *repositories.CardRepository
}

func NewCardService(cardRepo *repositories.CardRepository) *CardService {
	return &CardService{CardRepo: cardRepo}
}

// Генерация номера карты (зашифрованного)
func GenerateEncryptedCardNumber() (string, error) {
	plaintext := "4111111111111111"
	hash := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(hash[:]), nil
}

// Создание карты
func (s *CardService) CreateCard(accountID int, cvv string) (*models.Card, error) {
	encryptedPAN, err := GenerateEncryptedCardNumber()
	if err != nil {
		return nil, err
	}

	cvvHash, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	card := &models.Card{
		AccountID:    accountID,
		EncryptedPAN: encryptedPAN,
		CVVHash:      string(cvvHash),
		ExpiryMonth:  12,
		ExpiryYear:   2028,
	}

	err = s.CardRepo.CreateCard(card)
	if err != nil {
		return nil, err
	}

	return card, nil
}
