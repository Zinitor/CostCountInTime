package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"costcountintimeapp/models"
)

func TestCalculateHandler_Basic(t *testing.T) {
	tests := []struct {
		name         string
		salary       float64
		salaryPeriod string
		price        float64
		itemName     string
		wantStatus   int
	}{
		{"Hourly salary", 1000, "hour", 5000, "Наушники", http.StatusOK},
		{"Daily salary", 5000, "day", 15000, "Телефон", http.StatusOK},
		{"Monthly salary", 50000, "month", 100000, "Ноутбук", http.StatusOK},
		{"Empty item name", 50000, "month", 50000, "", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := models.CalculateRequest{
				Salary:       tt.salary,
				SalaryPeriod: tt.salaryPeriod,
				Price:        tt.price,
				ItemName:     tt.itemName,
			}

			jsonBody, _ := json.Marshal(body)
			req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			CalculateHandler(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}

			var response models.CalculateResponse
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				t.Errorf("failed to decode response: %v", err)
			}

			if response.Hours <= 0 {
				t.Error("hours should be positive")
			}

			if response.Formatted == "" {
				t.Error("formatted should not be empty")
			}
		})
	}
}

func TestCalculateHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/calculate", nil)
	rr := httptest.NewRecorder()
	CalculateHandler(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestCalculateHandler_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	CalculateHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestCalculateHandler_NegativeValues(t *testing.T) {
	tests := []struct {
		name   string
		body   models.CalculateRequest
	}{
		{"Negative salary", models.CalculateRequest{Salary: -50000, SalaryPeriod: "month", Price: 10000, ItemName: "Тест"}},
		{"Negative price", models.CalculateRequest{Salary: 50000, SalaryPeriod: "month", Price: -10000, ItemName: "Тест"}},
		{"Zero salary", models.CalculateRequest{Salary: 0, SalaryPeriod: "month", Price: 10000, ItemName: "Тест"}},
		{"Zero price", models.CalculateRequest{Salary: 50000, SalaryPeriod: "month", Price: 0, ItemName: "Тест"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			CalculateHandler(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
			}
		})
	}
}

func TestCalculateHandler_InvalidPeriod(t *testing.T) {
	body := models.CalculateRequest{
		Salary:       50000,
		SalaryPeriod: "invalid",
		Price:        10000,
		ItemName:     "Тест",
	}

	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	CalculateHandler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestCalculateHandler_Calculation(t *testing.T) {
	tests := []struct {
		name         string
		salary       float64
		salaryPeriod string
		price        float64
		wantHours    float64
	}{
		{"Hourly: 1000₽/hour, 5000₽ item", 1000, "hour", 5000, 5},
		{"Daily: 5000₽/day, 10000₽ item", 5000, "day", 10000, 16},
		{"Monthly: 50000₽/month, 100000₽ item", 50000, "month", 100000, 336},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := models.CalculateRequest{
				Salary:       tt.salary,
				SalaryPeriod: tt.salaryPeriod,
				Price:        tt.price,
				ItemName:     "Тест",
			}

			jsonBody, _ := json.Marshal(body)
			req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			CalculateHandler(rr, req)

			var response models.CalculateResponse
			json.NewDecoder(rr.Body).Decode(&response)

			// Проверяем с допустимой погрешностью
			diff := response.Hours - tt.wantHours
			if diff < -1 || diff > 1 {
				t.Errorf("hours = %v, want approximately %v", response.Hours, tt.wantHours)
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name   string
		hours  float64
		want   string
	}{
		{"Minutes", 0.5, "30 мин"},
		{"Less than 1 hour", 0.8, "48 мин"},
		{"1 hour", 1, "1 ч"},
		{"5 hours", 5, "5 ч"},
		{"1 day", 8, "1 дн"},
		{"1 day 4 hours", 12, "1 дн 4 ч"},
		{"5 days", 40, "5 дн"},
		{"10 days 4 hours", 84, "10 дн 4 ч"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatTime(tt.hours)
			if got != tt.want {
				t.Errorf("formatTime(%v) = %v, want %v", tt.hours, got, tt.want)
			}
		})
	}
}

func TestHoursInPeriod(t *testing.T) {
	tests := []struct {
		period string
		want   float64
	}{
		{"hour", 1},
		{"day", 8},
		{"month", 168},
	}

	for _, tt := range tests {
		t.Run(tt.period, func(t *testing.T) {
			got, ok := HoursInPeriod[tt.period]
			if !ok {
				t.Errorf("period %s not found in HoursInPeriod", tt.period)
			}
			if got != tt.want {
				t.Errorf("HoursInPeriod[%s] = %v, want %v", tt.period, got, tt.want)
			}
		})
	}
}

func TestCORS(t *testing.T) {
	handler := CORS(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Test OPTIONS
	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("OPTIONS returned wrong status: got %v want %v", status, http.StatusOK)
	}

	headers := rr.Header()
	if headers.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("missing CORS header")
	}
}
