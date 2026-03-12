package dto

import "costcountintimeapp/domain/valueobject"

// KarmaUpdateRequest представляет запрос на обновление кармы
type KarmaUpdateRequest struct {
	UserID      string  `json:"userId"`
	ItemName    string  `json:"item_name"`
	Price       float64 `json:"price"`
	HoursNeeded float64 `json:"hours_needed"`
	TagID       string  `json:"tag_id,omitempty"` // ID выбранного тега (опционально)
}

// KarmaUpdateResponse представляет ответ на обновление кармы
type KarmaUpdateResponse struct {
	Score           float64                `json:"score"`
	Level           valueobject.KarmaLevel `json:"level"`
	ScoreDelta      float64                `json:"score_delta"`
	Reasons         []string               `json:"reasons"`
	Stats           KarmaStats             `json:"stats"`
	Comment         string                 `json:"comment"`
	PurchaseComment string                 `json:"purchase_comment"`
}

// KarmaStateResponse представляет текущее состояние кармы
type KarmaStateResponse struct {
	Score float64                `json:"score"`
	Level valueobject.KarmaLevel `json:"level"`
	Info  valueobject.LevelInfo  `json:"info"`
	Stats KarmaStats             `json:"stats"`
}

// KarmaStats содержит статистику кармы
type KarmaStats struct {
	TotalSpent     float64 `json:"total_spent"`
	TotalHours     float64 `json:"total_hours"`
	AvgPurchase    float64 `json:"avg_purchase"`
	ImpulseCount   int     `json:"impulse_count"`
	BestStreak     int     `json:"best_streak"`
	CurrentStreak  int     `json:"current_streak"`
	PurchasesToday int     `json:"purchases_today"`
	SpentToday     float64 `json:"spent_today"`
}

// KarmaResetResponse представляет ответ на сброс кармы
type KarmaResetResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// KarmaDeleteRequest представляет запрос на удаление события кармы
type KarmaDeleteRequest struct {
	ItemName  string  `json:"itemName"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

// KarmaDeleteResponse представляет ответ на удаление события кармы
type KarmaDeleteResponse struct {
	Status       string `json:"status"`
	RowsAffected int64  `json:"rows_affected"`
}
