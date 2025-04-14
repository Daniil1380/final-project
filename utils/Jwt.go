package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

// GetUserIDFromContext извлекает ID пользователя из контекста запроса
func GetUserIDFromContext(r *http.Request) int {
	// Получаем значение из контекста
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		// Возвращаем 0, если ID пользователя не найден в контексте
		// Это можно обработать в вызывающем коде как ошибку аутентификации
		return 0
	}

	return userID
}

// GenerateJWT создает JWT-токен для пользователя
func GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}

// ValidateJWT проверяет и валидирует JWT-токен
func ValidateJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что используется ожидаемый алгоритм
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		// Возвращаем секретный ключ для проверки подписи
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Проверяем срок действия токена
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return 0, jwt.ErrTokenExpired
			}
		}

		// Извлекаем ID пользователя
		if userID, ok := claims["user_id"].(float64); ok {
			return int(userID), nil
		}

		return 0, jwt.ErrInvalidKey
	}

	return 0, jwt.ErrInvalidKey
}
