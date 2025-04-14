package services

import (
	"crypto/sha256"
	"encoding/hex"
	_ "errors"
	"final-project/models"
	"final-project/repositories"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
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

func GenerateCardNumber() (string, error) {
	// Префикс Visa
	prefix := "4"

	// Генерируем случайные цифры (всего 16 цифр, первая уже есть)
	rand.Seed(time.Now().UnixNano())
	digits := make([]int, 15)
	for i := 0; i < 15; i++ {
		digits[i] = rand.Intn(10)
	}

	// Конвертируем в строку
	cardNumberStr := prefix
	for _, digit := range digits {
		cardNumberStr += strconv.Itoa(digit)
	}

	// Проверяем по алгоритму Луна
	cardNumber := cardNumberStr[:15]
	sum := 0

	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(cardNumber[i]))

		if (len(cardNumber)-i)%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
	}

	// Вычисляем контрольную цифру
	checkDigit := (10 - (sum % 10)) % 10

	return cardNumberStr + strconv.Itoa(checkDigit), nil
}
