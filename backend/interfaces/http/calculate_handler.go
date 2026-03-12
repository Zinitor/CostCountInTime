package http

import (
	"encoding/json"
	"net/http"

	"costcountintimeapp/application/dto"
	"costcountintimeapp/application/service"
)

// CalculateHandler обрабатывает запросы на расчёт времени работы
type CalculateHandler struct {
	calculateService *service.CalculateService
}

// NewCalculateHandler создаёт новый CalculateHandler
func NewCalculateHandler(calculateService *service.CalculateService) *CalculateHandler {
	return &CalculateHandler{
		calculateService: calculateService,
	}
}

// ServeHTTP обрабатывает HTTP запрос
func (h *CalculateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.calculateService.Calculate(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
