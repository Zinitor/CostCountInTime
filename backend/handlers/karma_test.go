package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"costcountintimeapp/models"
)

// Инициализация тестовой БД
func setupTestDB(t *testing.T) func() {
	// Создаём временную БД
	dbPath := "test_karma.db"
	if err := InitKarmaDB(dbPath); err != nil {
		t.Fatalf("Failed to init test DB: %v", err)
	}

	// Функция очистки
	return func() {
		KarmaDB.Close()
		os.Remove(dbPath)
	}
}

func TestKarmaUpdateHandler_Basic(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	tests := []struct {
		name        string
		price       float64
		hoursNeeded float64
		itemName    string
		wantStatus  int
	}{
		{"Cheap item", 500, 0.5, "Салат", http.StatusOK},
		{"Normal item", 5000, 5, "Наушники", http.StatusOK},
		{"Expensive item", 50000, 50, "Телефон", http.StatusOK},
		{"Crazy item", 500000, 300, "Дача", http.StatusOK},
		{"Insane item", 5000000, 1000, "Дом", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := models.KarmaUpdateRequest{
				Price:       tt.price,
				HoursNeeded: tt.hoursNeeded,
				ItemName:    tt.itemName,
			}

			jsonBody, _ := json.Marshal(body)
			req := httptest.NewRequest(http.MethodPost, "/api/karma/update", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			KarmaUpdateHandler(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}

			var response models.KarmaUpdateResponse
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				t.Errorf("failed to decode response: %v", err)
			}

			if response.PurchaseComment == "" {
				t.Error("purchase_comment should not be empty")
			}
		})
	}
}

func TestKarmaUpdateHandler_InvalidMethod(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/karma/update", nil)
	rr := httptest.NewRecorder()
	KarmaUpdateHandler(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestKarmaUpdateHandler_InvalidBody(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodPost, "/api/karma/update", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	KarmaUpdateHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestKarmaStateHandler(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/karma/state", nil)
	rr := httptest.NewRecorder()
	KarmaStateHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if response["score"] == nil {
		t.Error("response should contain score")
	}

	if response["level"] == nil {
		t.Error("response should contain level")
	}
}

func TestKarmaResetHandler(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Сначала добавим данные
	body := models.KarmaUpdateRequest{
		Price:       5000,
		HoursNeeded: 10,
		ItemName:    "Тест",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/karma/update", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	KarmaUpdateHandler(rr, req)

	// Теперь сбрасываем
	req = httptest.NewRequest(http.MethodPost, "/api/karma/reset", nil)
	rr = httptest.NewRecorder()
	KarmaResetHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем что score = 0
	req = httptest.NewRequest(http.MethodGet, "/api/karma/state", nil)
	rr = httptest.NewRecorder()
	KarmaStateHandler(rr, req)

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	if score := response["score"].(float64); score != 0 {
		t.Errorf("score should be 0 after reset, got %v", score)
	}
}

func TestKarmaResetHandler_InvalidMethod(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/karma/reset", nil)
	rr := httptest.NewRecorder()
	KarmaResetHandler(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestKarmaDeleteEventHandler(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Добавим событие
	body := models.KarmaUpdateRequest{
		Price:       5000,
		HoursNeeded: 10,
		ItemName:    "Тест",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/karma/update", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	KarmaUpdateHandler(rr, req)

	// Получим timestamp
	var updateResp models.KarmaUpdateResponse
	json.NewDecoder(rr.Body).Decode(&updateResp)

	// Удаляем событие
	deleteBody := map[string]interface{}{
		"itemName":  "Тест",
		"price":     5000.0,
		"timestamp": getCurrentTimestamp(),
	}
	jsonDeleteBody, _ := json.Marshal(deleteBody)
	req = httptest.NewRequest(http.MethodPost, "/api/karma/delete_event", bytes.NewBuffer(jsonDeleteBody))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	KarmaDeleteEventHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestCalculateKarmaDelta(t *testing.T) {
	tests := []struct {
		name        string
		hoursNeeded float64
		wantDelta   float64
		wantComment bool
	}{
		{"Very cheap", 0.5, -2, true},
		{"Cheap", 5, 0, true},
		{"Medium", 20, 5, true},
		{"Expensive", 100, 15, true},
		{"Very expensive", 300, 25, true},
		{"Insane", 600, 35, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta, reasons, comment := calculateKarmaDelta(tt.hoursNeeded*1000, tt.hoursNeeded, "Тест")

			if delta != tt.wantDelta {
				t.Errorf("delta = %v, want %v", delta, tt.wantDelta)
			}

			if len(reasons) == 0 {
				t.Error("reasons should not be empty")
			}

			if tt.wantComment && comment == "" {
				t.Error("comment should not be empty")
			}
		})
	}
}

func TestGetKarmaLevel(t *testing.T) {
	tests := []struct {
		name      string
		score     float64
		wantLevel models.KarmaLevel
	}{
		{"Financial God", 10, models.FinancialGod},
		{"Wise Spender", 30, models.WiseSpender},
		{"Normal Person", 50, models.NormalPerson},
		{"Risky Buyer", 70, models.RiskyBuyer},
		{"Financial Degen", 90, models.FinancialDegen},
		{"Max score", 100, models.FinancialDegen},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := getKarmaLevel(tt.score)
			if level.Level != tt.wantLevel {
				t.Errorf("level = %v, want %v", level.Level, tt.wantLevel)
			}
		})
	}
}

func TestGetCurrentKarmaScore_Empty(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	score := getCurrentKarmaScore()
	if score != 0 {
		t.Errorf("score should be 0 for empty DB, got %v", score)
	}
}

func TestGetCurrentKarmaScore_Decay(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	// Добавим событие
	body := models.KarmaUpdateRequest{
		Price:       5000,
		HoursNeeded: 10,
		ItemName:    "Тест",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/karma/update", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	KarmaUpdateHandler(rr, req)

	score := getCurrentKarmaScore()
	if score <= 0 || score > 100 {
		t.Errorf("score should be between 0 and 100, got %v", score)
	}
}
