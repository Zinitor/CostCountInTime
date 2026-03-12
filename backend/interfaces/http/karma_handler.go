package http

import (
	"encoding/json"
	"net/http"

	"costcountintimeapp/application/dto"
	"costcountintimeapp/application/service"
)

// KarmaHandler обрабатывает запросы на управление кармой
type KarmaHandler struct {
	karmaService *service.KarmaService
}

// NewKarmaHandler создаёт новый KarmaHandler
func NewKarmaHandler(karmaService *service.KarmaService) *KarmaHandler {
	return &KarmaHandler{
		karmaService: karmaService,
	}
}

// UpdateHandler обрабатывает запрос на обновление кармы
func (h *KarmaHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.KarmaUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.karmaService.UpdateKarma(req)
	if err != nil {
		http.Error(w, "Failed to update karma", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// StateHandler обрабатывает запрос на получение состояния кармы
func (h *KarmaHandler) StateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp, err := h.karmaService.GetState()
	if err != nil {
		http.Error(w, "Failed to get karma state", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ResetHandler обрабатывает запрос на сброс кармы
func (h *KarmaHandler) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.karmaService.Reset(); err != nil {
		http.Error(w, "Failed to reset karma", http.StatusInternalServerError)
		return
	}

	resp := dto.KarmaResetResponse{
		Status:  "ok",
		Message: "Карма сброшена",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteEventHandler обрабатывает запрос на удаление события кармы
func (h *KarmaHandler) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.KarmaDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rowsAffected, err := h.karmaService.DeleteEvent(req.ItemName, req.Price, req.Timestamp)
	if err != nil {
		http.Error(w, "Failed to delete karma event", http.StatusInternalServerError)
		return
	}

	resp := dto.KarmaDeleteResponse{
		Status:       "ok",
		RowsAffected: rowsAffected,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
