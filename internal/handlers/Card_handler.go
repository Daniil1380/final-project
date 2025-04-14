package handlers

import (
	"encoding/json"
	"final-project/services"
	"final-project/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CardHandler struct {
	CardService *services.CardService
}

func NewCardHandler(cardService *services.CardService) *CardHandler {
	return &CardHandler{CardService: cardService}
}

// Создание виртуальной карты
func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID int    `json:"account_id"`
		CVV       string `json:"cvv"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	card, err := h.CardService.CreateCard(req.AccountID, req.CVV)
	if err != nil {
		http.Error(w, "Failed to create card", http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, card)
}

// Получение информации о карте (без CVV)
func (h *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	cardID, _ := strconv.Atoi(mux.Vars(r)["id"])

	card, err := h.CardService.CardRepo.GetCardByID(cardID)
	if err != nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	utils.RespondJSON(w, http.StatusOK, card)
}
