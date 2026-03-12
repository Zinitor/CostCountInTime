package dto

// CalculateRequest представляет запрос на расчёт времени работы
type CalculateRequest struct {
	Salary       float64 `json:"salary"`
	SalaryPeriod string  `json:"salaryPeriod"`
	Price        float64 `json:"price"`
	ItemName     string  `json:"itemName"`
	TagID        string  `json:"tagId,omitempty"` // ID выбранного тега (опционально)
}

// CalculateResponse представляет ответ на расчёт времени работы
type CalculateResponse struct {
	Hours     float64 `json:"hours"`
	Days      float64 `json:"days"`
	Formatted string  `json:"formatted"`
	ItemName  string  `json:"itemName"`
	Price     float64 `json:"price"`
	Comment   string  `json:"comment"`
	TagID     string  `json:"tagId,omitempty"`   // ID выбранного тега
	TagName   string  `json:"tagName,omitempty"` // Название тега
}
