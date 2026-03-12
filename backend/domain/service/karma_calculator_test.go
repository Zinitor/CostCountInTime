package service_test

import (
	"testing"
	"time"

	"costcountintimeapp/domain/entity"
	"costcountintimeapp/domain/repository"
	"costcountintimeapp/domain/service"
	"costcountintimeapp/domain/valueobject"
	"costcountintimeapp/infrastructure/persistence"
)

func setupTestTagRepo() repository.TagRepository {
	return persistence.NewInMemoryTagRepository()
}

func TestKarmaCalculatorService_Calculate(t *testing.T) {
	tagRepo := setupTestTagRepo()
	svc := service.NewKarmaCalculatorService(tagRepo)
	stats := service.KarmaStats{}

	tests := []struct {
		name        string
		hoursNeeded float64
		wantDelta   float64
	}{
		{"Very cheap (< 2 hours)", 0.5, -2},
		{"Cheap (2-40 hours)", 5, 0},
		{"Medium (40-200 hours)", 100, 5},
		{"Expensive (200-500 hours)", 300, 15},
		{"Very expensive (500+ hours)", 600, 25},
		{"Insane (500+ hours default)", 2000, 35},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.Calculate(tt.hoursNeeded*1000, tt.hoursNeeded, "Тест", nil, stats, "")

			if result.Delta != tt.wantDelta {
				t.Errorf("delta = %v, want %v", result.Delta, tt.wantDelta)
			}

			if len(result.Reasons) == 0 {
				t.Error("reasons should not be empty")
			}

			if result.PurchaseComment == "" {
				t.Error("comment should not be empty")
			}
		})
	}
}

func TestKarmaCalculatorService_GetLevelByScore(t *testing.T) {
	tagRepo := setupTestTagRepo()
	svc := service.NewKarmaCalculatorService(tagRepo)

	tests := []struct {
		name      string
		score     float64
		wantLevel string
	}{
		{"Financial God", 10, "financial_god"},
		{"Wise Spender", 30, "wise_spender"},
		{"Normal Person", 50, "normal_person"},
		{"Risky Buyer", 70, "risky_buyer"},
		{"Financial Degen", 90, "financial_degen"},
		{"Max score", 100, "financial_degen"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := svc.GetLevelByScore(tt.score)
			if string(level.Level) != tt.wantLevel {
				t.Errorf("level = %v, want %v", level.Level, tt.wantLevel)
			}
		})
	}
}

func TestKarmaCalculatorService_GetCommentByLevel(t *testing.T) {
	tagRepo := setupTestTagRepo()
	svc := service.NewKarmaCalculatorService(tagRepo)

	tests := []struct {
		name    string
		level   string
		delta   float64
		reasons []string
	}{
		{"Financial God comment", "financial_god", 0, []string{}},
		{"Financial Degen comment", "financial_degen", 20, []string{"Reason 1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment := svc.GetCommentByLevel(
				valueobject.KarmaLevel(tt.level),
				tt.delta,
				tt.reasons,
			)

			if comment == "" {
				t.Error("comment should not be empty")
			}
		})
	}
}

func TestKarmaCalculatorService_ImpulseBuying(t *testing.T) {
	tagRepo := setupTestTagRepo()
	svc := service.NewKarmaCalculatorService(tagRepo)
	stats := service.KarmaStats{}

	// Создаём последнее событие (покупка была только что)
	now := time.Now().Unix()
	lastEvent := &entity.KarmaEvent{
		Timestamp: now - 1800, // 30 минут назад (в пределах часа)
	}

	result := svc.Calculate(5000, 5, "Тест", lastEvent, stats, "")

	// Проверяем что расчёт работает
	if len(result.Reasons) == 0 {
		t.Error("calculation should return non-empty reasons")
	}
}
