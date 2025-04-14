package main

import (
	"database/sql"
	"final-project/internal/handlers"
	"final-project/internal/middleware"
	"final-project/repositories"
	"final-project/scheduler"
	"final-project/services"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Подключение к БД
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Репозитории
	userRepo := repositories.NewUserRepository(db)
	accountRepo := repositories.NewAccountRepository(db)
	cardRepo := repositories.NewCardRepository(db)
	loanRepo := repositories.NewLoanRepository(db) // добавляем кредитный репозиторий

	// Сервисы
	userService := services.NewUserService(userRepo)
	accountService := services.NewAccountService(accountRepo)
	cardService := services.NewCardService(cardRepo)
	loanService := services.NewLoanService(loanRepo, accountRepo)

	// Обработчики
	authHandler := handlers.NewAuthHandler(userService)
	accountHandler := handlers.NewAccountHandler(accountService)
	cardHandler := handlers.NewCardHandler(cardService)
	// Предположим, позже будут реализованы обработчики для кредитов

	// Запуск шедулера для автоматического списания платежей по кредитам
	scheduler.StartCreditPaymentScheduler(loanService)

	// Маршрутизация
	r := mux.NewRouter()

	// Публичные маршруты
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Защищённые маршруты
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	api.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	api.HandleFunc("/accounts/{id}/deposit", accountHandler.Deposit).Methods("POST")
	api.HandleFunc("/accounts/{id}/withdraw", accountHandler.Withdraw).Methods("POST")
	api.HandleFunc("/cards", cardHandler.CreateCard).Methods("POST")
	api.HandleFunc("/cards/{id}", cardHandler.GetCard).Methods("GET")
	// Эндпоинты для кредитов можно добавить здесь

	// Запуск сервера
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
