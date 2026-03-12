package entity

import "costcountintimeapp/domain/valueobject"

// KarmaEvent представляет событие изменения кармы
type KarmaEvent struct {
	ID          int64
	UserID      string
	ItemName    string
	Price       float64
	HoursNeeded float64
	ScoreDelta  float64
	Reason      string
	Timestamp   int64
}

// NewKarmaEvent создаёт новое событие кармы
func NewKarmaEvent(userID, itemName string, price, hoursNeeded, scoreDelta float64, reason string, timestamp int64) *KarmaEvent {
	return &KarmaEvent{
		UserID:      userID,
		ItemName:    itemName,
		Price:       price,
		HoursNeeded: hoursNeeded,
		ScoreDelta:  scoreDelta,
		Reason:      reason,
		Timestamp:   timestamp,
	}
}

// GetLevel возвращает уровень кармы на основе события
func (e *KarmaEvent) GetLevel() valueobject.KarmaLevel {
	for _, level := range valueobject.KarmaLevels {
		if e.ScoreDelta >= level.MinScore && e.ScoreDelta < level.MaxScore {
			return level.Level
		}
	}
	return valueobject.FinancialDegen
}
