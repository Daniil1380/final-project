package main

import (
	"database/sql"
	"final-project/internal/handlers"
	"final-project/internal/middleware"
	"final-project/repositories"
	"final-project/scheduler"
	"final-project/services"
	"final-project/utils"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "PAASSSSSSSS"
	dbname   = "finance"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// Инициализация логгера
	utils.InitLogger()
	utils.Logger.Info("Starting application...")

	// Подключение к БД
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		utils.Logger.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Репозитории
	userRepo := repositories.NewUserRepository(db)
	accountRepo := repositories.NewAccountRepository(db)
	cardRepo := repositories.NewCardRepository(db)
	loanRepo := repositories.NewLoanRepository(db)
	paymentScheduleRepo := repositories.NewPaymentScheduleRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Сервисы
	userService := services.NewUserService(userRepo)
	accountService := services.NewAccountService(accountRepo, transactionRepo)
	cardService := services.NewCardService(cardRepo)
	loanService := services.NewLoanService(loanRepo, accountRepo, paymentScheduleRepo)

	// Обработчики
	authHandler := handlers.NewAuthHandler(userService)
	accountHandler := handlers.NewAccountHandler(accountService)
	cardHandler := handlers.NewCardHandler(cardService)
	loanHandler := handlers.NewLoanHandler(loanService)

	// Запуск шедулера для автоматического списания платежей по кредитам
	scheduler.StartCreditPaymentScheduler(loanService)

	// Маршрутизация
	r := mux.NewRouter()

	// Публичные эндпоинты
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Защищённые эндпоинты (JWT через middleware)
	api := r.PathPrefix("/").Subrouter()
	api.Use(middleware.AuthMiddleware)

	// Счета
	api.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	api.HandleFunc("/accounts/{id}/deposit", accountHandler.Deposit).Methods("POST")
	api.HandleFunc("/accounts/{id}/withdraw", accountHandler.Withdraw).Methods("POST")

	// Переводы
	api.HandleFunc("/transfer", accountHandler.Transfer).Methods("POST")

	// Карты
	api.HandleFunc("/cards", cardHandler.CreateCard).Methods("POST")
	api.HandleFunc("/cards/{id}", cardHandler.GetCard).Methods("GET")

	// Аналитика
	api.HandleFunc("/analytics", accountHandler.GetAccountAnalytics).Methods("GET")

	// Кредиты
	api.HandleFunc("/credits", loanHandler.CreateLoan).Methods("POST")
	api.HandleFunc("/credits/{id}", loanHandler.GetLoan).Methods("GET")
	api.HandleFunc("/credits/{creditId}/schedule", loanHandler.GetPaymentSchedule).Methods("GET")
	api.HandleFunc("/credits/process-payments", loanHandler.ProcessPayments).Methods("POST")

	utils.Logger.Info("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
