package handlers

import (
	"encoding/json"
	"math"
	"net/http"

	"costcountintimeapp/models"
	"time"
)

// KarmaUpdateHandlerSimple - in-memory version for Vercel
func KarmaUpdateHandlerSimple(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.KarmaUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	delta, reasons, purchaseComment := calculateKarmaDelta(req.Price, req.HoursNeeded, req.ItemName)

	event := &models.KarmaEvent{
		UserID:      "default",
		ItemName:    req.ItemName,
		Price:       req.Price,
		HoursNeeded: req.HoursNeeded,
		ScoreDelta:  delta,
		Reason:      joinStrings(reasons, "; "),
		Timestamp:   time.Now().Unix(),
	}

	if err := saveKarmaEventSimple(event); err != nil {
		http.Error(w, "Failed to save karma event", http.StatusInternalServerError)
		return
	}

	newScore := getCurrentKarmaScoreSimple()
	level := getKarmaLevel(newScore)
	stats, _ := getKarmaStatsSimple()
	comment := getKarmaComment(level.Level, delta, reasons)

	response := models.KarmaUpdateResponse{
		Score:           math.Round(newScore*100) / 100,
		Level:           level.Level,
		ScoreDelta:      delta,
		Reasons:         reasons,
		Stats:           stats,
		Comment:         comment,
		PurchaseComment: purchaseComment,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// KarmaStateHandlerSimple - in-memory version for Vercel
func KarmaStateHandlerSimple(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	score := getCurrentKarmaScoreSimple()
	level := getKarmaLevel(score)
	stats, _ := getKarmaStatsSimple()

	response := map[string]interface{}{
		"score": score,
		"level": level.Level,
		"info":  level,
		"stats": stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// KarmaResetHandlerSimple - in-memory version for Vercel
func KarmaResetHandlerSimple(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := resetKarmaSimple(); err != nil {
		http.Error(w, "Failed to reset karma", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":  "ok",
		"message": "Карма сброшена",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// KarmaDeleteEventHandlerSimple - in-memory version for Vercel
func KarmaDeleteEventHandlerSimple(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ItemName  string  `json:"itemName"`
		Price     float64 `json:"price"`
		Timestamp int64   `json:"timestamp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	deleted, err := deleteKarmaEventSimple(req.ItemName, req.Price, req.Timestamp)
	if err != nil {
		http.Error(w, "Failed to delete karma event", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":        "ok",
		"rows_affected": deleted,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
