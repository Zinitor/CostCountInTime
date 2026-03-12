package service_test

import (
	"testing"

	"costcountintimeapp/application/dto"
	"costcountintimeapp/application/service"
	"costcountintimeapp/infrastructure/persistence"
)

func setupTestTagRepo() *persistence.InMemoryTagRepository {
	return persistence.NewInMemoryTagRepository()
}

func TestCalculateService_Basic(t *testing.T) {
	tagRepo := setupTestTagRepo()
	svc := service.NewCalculateService(tagRepo)

	tests := []struct {
		name         string
		salary       float64
		salaryPeriod string
		price        float64
		itemName     string
		wantErr      bool
	}{
		{"Hourly salary", 1000, "hour", 5000, "Наушники", false},
		{"Daily salary", 5000, "day", 15000, "Телефон", false},
		{"Monthly salary", 50000, "month", 100000, "Ноутбук", false},
		{"Empty item name", 50000, "month", 50000, "", false},
		{"Negative salary", -50000, "month", 10000, "Тест", true},
		{"Negative price", 50000, "month", -10000, "Тест", true},
		{"Zero salary", 0, "month", 10000, "Тест", true},
		{"Zero price", 50000, "month", 0, "Тест", true},
		{"Invalid period", 50000, "invalid", 10000, "Тест", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := dto.CalculateRequest{
				Salary:       tt.salary,
				SalaryPeriod: tt.salaryPeriod,
				Price:        tt.price,
				ItemName:     tt.itemName,
			}

			resp, err := svc.Calculate(req)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.wantErr {
				if resp.Hours <= 0 {
					t.Error("hours should be positive")
				}
				if resp.Formatted == "" {
					t.Error("formatted should not be empty")
				}
				if resp.Comment == "" {
					t.Error("comment should not be empty")
				}
			}
		})
	}
}

func TestCalculateService_Calculation(t *testing.T) {
	tagRepo := setupTestTagRepo()
	svc := service.NewCalculateService(tagRepo)

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
			req := dto.CalculateRequest{
				Salary:       tt.salary,
				SalaryPeriod: tt.salaryPeriod,
				Price:        tt.price,
				ItemName:     "Тест",
			}

			resp, _ := svc.Calculate(req)

			// Проверяем с допустимой погрешностью
			diff := resp.Hours - tt.wantHours
			if diff < -1 || diff > 1 {
				t.Errorf("hours = %v, want approximately %v", resp.Hours, tt.wantHours)
			}
		})
	}
}
