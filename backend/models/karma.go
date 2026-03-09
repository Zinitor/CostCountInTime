package models

type KarmaLevel string

const (
	FinancialGod   KarmaLevel = "financial_god"
	WiseSpender    KarmaLevel = "wise_spender"
	NormalPerson   KarmaLevel = "normal_person"
	RiskyBuyer     KarmaLevel = "risky_buyer"
	FinancialDegen KarmaLevel = "financial_degen"
)

type KarmaState struct {
	UserID      string     `json:"user_id"`
	Score       float64    `json:"score"`
	Level       KarmaLevel `json:"level"`
	LastUpdated int64      `json:"last_updated"`
}

type KarmaEvent struct {
	ID          int64   `json:"id"`
	UserID      string  `json:"user_id"`
	ItemName    string  `json:"item_name"`
	Price       float64 `json:"price"`
	HoursNeeded float64 `json:"hours_needed"`
	ScoreDelta  float64 `json:"score_delta"`
	Reason      string  `json:"reason"`
	Timestamp   int64   `json:"timestamp"`
}

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

type KarmaUpdateRequest struct {
	Salary       float64 `json:"salary"`
	SalaryPeriod string  `json:"salary_period"`
	Price        float64 `json:"price"`
	ItemName     string  `json:"item_name"`
	HoursNeeded  float64 `json:"hours_needed"`
}

type KarmaUpdateResponse struct {
	Score           float64    `json:"score"`
	Level           KarmaLevel `json:"level"`
	ScoreDelta      float64    `json:"score_delta"`
	Reasons         []string   `json:"reasons"`
	Stats           KarmaStats `json:"stats"`
	Comment         string     `json:"comment"`          // Комментарий уровня кармы
	PurchaseComment string     `json:"purchase_comment"` // Комментарий покупки
}

type LevelInfo struct {
	Level       KarmaLevel
	MinScore    float64
	MaxScore    float64
	DisplayName string
	Icon        string
	Color       string
}
