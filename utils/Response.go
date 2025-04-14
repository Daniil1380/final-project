package utils

import (
	"encoding/json"
	"net/http"
)

// RespondJSON отправляет JSON-ответ клиенту
func RespondJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload != nil {
		err := json.NewEncoder(w).Encode(payload)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
