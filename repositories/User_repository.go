package repositories

import (
	"database/sql"
	"errors"
	"final-project/models"
)

type UserRepository struct {
	DB *sql.DB
}

// Создание нового экземпляра UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Создание пользователя
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, email, password, created_at, updated_at) 
                  VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`
	err := r.DB.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

// Получение пользователя по email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users WHERE email = $1`
	row := r.DB.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
