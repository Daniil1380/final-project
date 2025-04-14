package services

import (
	"errors"
	"final-project/models"
	"final-project/repositories"
	"golang.org/x/crypto/bcrypt"
)

// UserService управляет логикой работы с пользователями
type UserService struct {
	UserRepo *repositories.UserRepository
}

// Создание нового экземпляра UserService
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

// Регистрация нового пользователя
func (s *UserService) RegisterUser(username, email, password string) (*models.User, error) {
	// Проверяем, существует ли пользователь
	existingUser, _ := s.UserRepo.GetUserByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	// Создаём пользователя в БД
	err = s.UserRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Аутентификация пользователя
func (s *UserService) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}
