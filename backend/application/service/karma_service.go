package service

import (
	"math"
	"time"

	"costcountintimeapp/application/dto"
	"costcountintimeapp/domain/entity"
	"costcountintimeapp/domain/repository"
	domainService "costcountintimeapp/domain/service"
)

// KarmaService отвечает за управление кармой
type KarmaService struct {
	karmaRepo     repository.KarmaRepository
	tagRepo       repository.TagRepository
	calculator    *domainService.KarmaCalculatorService
	defaultUserID string
}

// NewKarmaService создаёт новый KarmaService
func NewKarmaService(karmaRepo repository.KarmaRepository, tagRepo repository.TagRepository) *KarmaService {
	return &KarmaService{
		karmaRepo:     karmaRepo,
		tagRepo:       tagRepo,
		calculator:    domainService.NewKarmaCalculatorService(tagRepo),
		defaultUserID: "default",
	}
}

// UpdateKarma обновляет карму пользователя
func (s *KarmaService) UpdateKarma(req dto.KarmaUpdateRequest) (*dto.KarmaUpdateResponse, error) {
	// Получаем последнее событие
	lastEvent, _ := s.karmaRepo.GetLast(s.defaultUserID)

	// Получаем статистику
	stats := s.getKarmaStats()

	// Рассчитываем изменение кармы с учётом тега
	result := s.calculator.Calculate(req.Price, req.HoursNeeded, req.ItemName, lastEvent, stats, req.TagID)

	// Создаём и сохраняем событие
	event := entity.NewKarmaEvent(
		s.defaultUserID,
		req.ItemName,
		req.Price,
		req.HoursNeeded,
		result.Delta,
		s.joinStrings(result.Reasons, "; "),
		time.Now().Unix(),
	)

	if err := s.karmaRepo.Save(event); err != nil {
		return nil, err
	}

	// Получаем новый счёт
	newScore := s.getCurrentKarmaScore()
	level := s.calculator.GetLevelByScore(newScore)
	comment := s.calculator.GetCommentByLevel(level.Level, result.Delta, result.Reasons)

	return &dto.KarmaUpdateResponse{
		Score:           math.Round(newScore*100) / 100,
		Level:           level.Level,
		ScoreDelta:      result.Delta,
		Reasons:         result.Reasons,
		Stats:           s.convertStats(stats),
		Comment:         comment,
		PurchaseComment: result.PurchaseComment,
	}, nil
}

// GetState возвращает текущее состояние кармы
func (s *KarmaService) GetState() (*dto.KarmaStateResponse, error) {
	score := s.getCurrentKarmaScore()
	level := s.calculator.GetLevelByScore(score)
	stats := s.getKarmaStats()

	return &dto.KarmaStateResponse{
		Score: score,
		Level: level.Level,
		Info:  level,
		Stats: s.convertStats(stats),
	}, nil
}

// Reset сбрасывает карму пользователя
func (s *KarmaService) Reset() error {
	return s.karmaRepo.DeleteAllByUser(s.defaultUserID)
}

// DeleteEvent удаляет событие кармы
func (s *KarmaService) DeleteEvent(itemName string, price float64, timestamp int64) (int64, error) {
	return s.karmaRepo.DeleteByItemAndTimestamp(s.defaultUserID, itemName, price, timestamp)
}

func (s *KarmaService) getCurrentKarmaScore() float64 {
	score, err := s.karmaRepo.GetTotalScore(s.defaultUserID)
	if err != nil {
		return 0
	}

	// Decay: уменьшаем на 5% за каждый день
	events, err := s.karmaRepo.GetAllByUser(s.defaultUserID)
	if err != nil || len(events) == 0 {
		return score
	}

	// Находим первое событие
	var firstTimestamp int64
	for _, event := range events {
		if firstTimestamp == 0 || event.Timestamp < firstTimestamp {
			firstTimestamp = event.Timestamp
		}
	}

	if firstTimestamp > 0 {
		daysPassed := float64(time.Now().Unix()-firstTimestamp) / 86400
		decay := math.Pow(0.95, daysPassed)
		score *= decay
	}

	// Ограничиваем 0-100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

func (s *KarmaService) getKarmaStats() domainService.KarmaStats {
	events, _ := s.karmaRepo.GetAllByUser(s.defaultUserID)

	stats := domainService.KarmaStats{}
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()

	for _, event := range events {
		stats.TotalSpent += event.Price
		stats.TotalHours += event.HoursNeeded

		if event.Timestamp >= todayStart {
			stats.PurchasesToday++
			stats.SpentToday += event.Price
		}
	}

	if len(events) > 0 {
		stats.AvgPurchase = stats.TotalSpent / float64(len(events))
	}

	return stats
}

func (s *KarmaService) convertStats(stats domainService.KarmaStats) dto.KarmaStats {
	return dto.KarmaStats{
		TotalSpent:     stats.TotalSpent,
		TotalHours:     stats.TotalHours,
		AvgPurchase:    stats.AvgPurchase,
		ImpulseCount:   stats.ImpulseCount,
		BestStreak:     stats.BestStreak,
		CurrentStreak:  stats.CurrentStreak,
		PurchasesToday: stats.PurchasesToday,
		SpentToday:     stats.SpentToday,
	}
}

func (s *KarmaService) joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
