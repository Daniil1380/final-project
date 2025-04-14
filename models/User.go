package models

import (
	"errors"
	"strings"
	"time"
)

// User представляет пользователя системы
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Хешированный пароль
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	// Проверка email
	if !strings.Contains(u.Email, "@") {
		return errors.New("invalid email format")
	}

	// Проверка имени пользователя
	if len(u.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	return nil
}
