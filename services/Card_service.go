package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"final-project/models"
	"final-project/repositories"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/openpgp"
	"io"
	"math/big"
	"os"
	"strconv"
	"strings"
)

type CardService struct {
	CardRepo *repositories.CardRepository
}

func NewCardService(cardRepo *repositories.CardRepository) *CardService {
	return &CardService{CardRepo: cardRepo}
}

// GenerateLuhnCardNumber генерирует номер карты по алгоритму Луна.
// Генерируем 15 случайных цифр с префиксом "4", затем вычисляем контрольную цифру.
func GenerateLuhnCardNumber() (string, error) {
	prefix := "4"
	// Генерируем 14 случайных цифр (итого 15 цифр без контрольной)
	numDigits := 14
	number := prefix
	for i := 0; i < numDigits; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		number += n.String()
	}
	// Вычисляем контрольную цифру по алгоритму Луна
	sum := 0
	reverseDigits := reverseString(number)
	for i, r := range reverseDigits {
		digit, _ := strconv.Atoi(string(r))
		// Если позиция (начиная с 0) четная – удваиваем цифру
		if i%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	checkDigit := (10 - (sum % 10)) % 10
	fullNumber := number + strconv.Itoa(checkDigit)
	return fullNumber, nil
}

// Вспомогательная функция для разворота строки.
func reverseString(s string) string {
	runes := []rune(s)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-i-1] = runes[n-i-1], runes[i]
	}
	return string(runes)
}

// PGPEncryptAndHMAC шифрует входящую строку (номер карты) с использованием PGP и вычисляет HMAC.
func PGPEncryptAndHMAC(cardNumber string) (encrypted string, hmacSignature string, err error) {
	// Получаем PGP-публичный ключ из переменных окружения
	pubKeyArmor := os.Getenv("PGP_PUBLIC_KEY")
	if pubKeyArmor == "" {
		return "", "", errors.New("PGP_PUBLIC_KEY не установлен")
	}
	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(pubKeyArmor))
	if err != nil {
		return "", "", fmt.Errorf("ошибка чтения PGP ключа: %v", err)
	}

	// Шифрование поля cardNumber
	var encryptedBuf bytes.Buffer
	encryptWriter, err := openpgp.Encrypt(&encryptedBuf, entityList, nil, nil, nil)
	if err != nil {
		return "", "", fmt.Errorf("ошибка шифрования: %v", err)
	}
	_, err = io.WriteString(encryptWriter, cardNumber)
	if err != nil {
		return "", "", err
	}
	err = encryptWriter.Close()
	if err != nil {
		return "", "", err
	}
	encrypted = encryptedBuf.String()

	// Получаем секрет для HMAC из переменных окружения
	hmacSecret := os.Getenv("HMAC_SECRET")
	if hmacSecret == "" {
		return "", "", errors.New("HMAC_SECRET не установлен")
	}
	h := hmac.New(sha256.New, []byte(hmacSecret))
	h.Write([]byte(encrypted))
	hmacSignature = hex.EncodeToString(h.Sum(nil))

	return encrypted, hmacSignature, nil
}

// CreateCard создает виртуальную карту, генерируя корректный номер по Луна, шифруя его, а также хешируя CVV.
func (s *CardService) CreateCard(accountID int, cvv string) (*models.Card, error) {
	// Генерируем номер карты по алгоритму Луна
	cardNumber, err := GenerateLuhnCardNumber()
	if err != nil {
		return nil, fmt.Errorf("не удалось сгенерировать номер карты: %v", err)
	}

	// Шифруем номер карты с использованием PGP и вычисляем HMAC подпись.
	encryptedPAN, hmacSignature, err := PGPEncryptAndHMAC(cardNumber)
	if err != nil {
		return nil, fmt.Errorf("ошибка шифрования карты: %v", err)
	}

	// Хешируем CVV с использованием bcrypt
	cvvHash, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Заполняем модель карты.
	// Для хранения HMAC подписи можно расширить модель (например, добавить поле DataHMAC).
	card := &models.Card{
		AccountID:    accountID,
		EncryptedPAN: encryptedPAN,
		CVVHash:      string(cvvHash),
		ExpiryMonth:  12,
		ExpiryYear:   2028,
		// Если требуется сохранить HMAC – можно добавить: DataHMAC: hmacSignature,
	}

	// Можно логировать для отладки (не выводить реальные данные в продакшене)
	fmt.Printf("Создана карта: зашифрованный PAN=%s, HMAC=%s\n", encryptedPAN, hmacSignature)

	err = s.CardRepo.CreateCard(card)
	if err != nil {
		return nil, err
	}

	return card, nil
}
