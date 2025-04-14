package handlers

import (
	"encoding/json"
	"final-project/services"
	"final-project/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AccountHandler struct {
	AccountService *services.AccountService
}

func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{AccountService: accountService}
}

// Создание банковского счета
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r)
	var req struct {
		Currency string `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	account, err := h.AccountService.CreateAccount(userID, req.Currency)
	if err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, account)
}

// Пополнение счета
func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	accountID, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := h.AccountService.Deposit(accountID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Deposit successful"})
}

// Списание средств
func (h *AccountHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	accountID, _ := strconv.Atoi(mux.Vars(r)["id"])
	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := h.AccountService.Withdraw(accountID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Withdrawal successful"})
}

func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r)

	vars := mux.Vars(r)
	fromAccountID, err := strconv.Atoi(vars["fromAccountID"])
	if err != nil {
		http.Error(w, "Invalid from account ID", http.StatusBadRequest)
		return
	}

	toAccountID, err := strconv.Atoi(vars["toAccountID"])
	if err != nil {
		http.Error(w, "Invalid to account ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.AccountService.TransferFunds(fromAccountID, toAccountID, req.Amount, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Transfer successful"})
}
