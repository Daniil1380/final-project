package handlers

import (
	"encoding/json"
	"final-project/services"
	"final-project/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type LoanHandler struct {
	LoanService *services.LoanService
}

// Создание нового LoanHandler
func NewLoanHandler(loanService *services.LoanService) *LoanHandler {
	return &LoanHandler{LoanService: loanService}
}

// Создание нового кредита
func (h *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r) // Получаем ID пользователя из контекста

	var req struct {
		AccountID int     `json:"account_id"`
		Amount    float64 `json:"amount"`
		Term      int     `json:"term"` // Срок кредита в месяцах
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	loan, err := h.LoanService.CreateLoan(userID, req.AccountID, req.Amount, req.Term)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, loan)
}

// Получение информации о кредите
func (h *LoanHandler) GetLoan(w http.ResponseWriter, r *http.Request) {
	loanID, err := strconv.Atoi(mux.Vars(r)["id"]) // ID кредита из URL
	if err != nil {
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	loan, err := h.LoanService.LoanRepo.GetLoanByID(loanID) // Предполагаем, что метод существует
	if err != nil {
		http.Error(w, "Loan not found", http.StatusNotFound)
		return
	}

	utils.RespondJSON(w, http.StatusOK, loan)
}

// Обработка ежемесячных платежей по всем активным кредитам
func (h *LoanHandler) ProcessPayments(w http.ResponseWriter, r *http.Request) {
	err := h.LoanService.ProcessCreditPayments()
	if err != nil {
		http.Error(w, "Failed to process payments", http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Payments processed successfully"})
}
