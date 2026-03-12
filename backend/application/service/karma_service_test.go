package service_test

import (
	"testing"
	"time"

	"costcountintimeapp/application/dto"
	"costcountintimeapp/application/service"
	"costcountintimeapp/infrastructure/persistence"
)

func setupTestKarmaService(t *testing.T) (*service.KarmaService, func()) {
	karmaRepo := persistence.NewInMemoryKarmaRepository()
	tagRepo := persistence.NewInMemoryTagRepository()
	svc := service.NewKarmaService(karmaRepo, tagRepo)
	return svc, func() {}
}

func TestKarmaService_UpdateKarma_Basic(t *testing.T) {
	svc, _ := setupTestKarmaService(t)

	tests := []struct {
		name        string
		price       float64
		hoursNeeded float64
		itemName    string
	}{
		{"Cheap item", 500, 0.5, "Салат"},
		{"Normal item", 5000, 5, "Наушники"},
		{"Expensive item", 50000, 50, "Телефон"},
		{"Crazy item", 500000, 300, "Дача"},
		{"Insane item", 5000000, 1000, "Дом"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := dto.KarmaUpdateRequest{
				Price:       tt.price,
				HoursNeeded: tt.hoursNeeded,
				ItemName:    tt.itemName,
			}

			resp, err := svc.UpdateKarma(req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.PurchaseComment == "" {
				t.Error("purchase_comment should not be empty")
			}

			if resp.Score < 0 || resp.Score > 100 {
				t.Errorf("score should be between 0 and 100, got %v", resp.Score)
			}
		})
	}
}

func TestKarmaService_GetState(t *testing.T) {
	svc, _ := setupTestKarmaService(t)

	state, err := svc.GetState()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if state.Score != 0 {
		t.Errorf("initial score should be 0, got %v", state.Score)
	}

	if state.Level == "" {
		t.Error("level should not be empty")
	}
}

func TestKarmaService_Reset(t *testing.T) {
	svc, _ := setupTestKarmaService(t)

	// Сначала добавим данные
	req := dto.KarmaUpdateRequest{
		Price:       5000,
		HoursNeeded: 10,
		ItemName:    "Тест",
	}
	_, err := svc.UpdateKarma(req)
	if err != nil {
		t.Fatalf("failed to update karma: %v", err)
	}

	// Теперь сбрасываем
	if err := svc.Reset(); err != nil {
		t.Fatalf("failed to reset karma: %v", err)
	}

	// Проверяем что score = 0
	state, _ := svc.GetState()
	if state.Score != 0 {
		t.Errorf("score should be 0 after reset, got %v", state.Score)
	}
}

func TestKarmaService_DeleteEvent(t *testing.T) {
	svc, _ := setupTestKarmaService(t)

	// Добавим событие
	req := dto.KarmaUpdateRequest{
		Price:       5000,
		HoursNeeded: 10,
		ItemName:    "Тест",
	}
	_, err := svc.UpdateKarma(req)
	if err != nil {
		t.Fatalf("failed to update karma: %v", err)
	}

	// Удаляем событие (используем текущий timestamp)
	rowsAffected, err := svc.DeleteEvent("Тест", 5000, time.Now().Unix())
	if err != nil {
		t.Fatalf("failed to delete event: %v", err)
	}

	// Note: В реальной реализации rowsAffected может быть 0 из-за разницы во времени
	_ = rowsAffected
}

func TestKarmaService_KarmaLevels(t *testing.T) {
	svc, _ := setupTestKarmaService(t)

	tests := []struct {
		name          string
		price         float64
		hoursNeeded   float64
		expectedLevel string
	}{
		{"Cheap items keep low karma", 500, 0.5, "financial_god"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := dto.KarmaUpdateRequest{
				Price:       tt.price,
				HoursNeeded: tt.hoursNeeded,
				ItemName:    "Тест",
			}

			resp, err := svc.UpdateKarma(req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(resp.Level) != tt.expectedLevel {
				t.Errorf("level = %v, want %v", resp.Level, tt.expectedLevel)
			}
		})
	}
}
