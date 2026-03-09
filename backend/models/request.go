package models

type CalculateRequest struct {
	Salary       float64 `json:"salary"`
	SalaryPeriod string  `json:"salaryPeriod"` // "hour", "day", "month"
	Price        float64 `json:"price"`
	ItemName     string  `json:"itemName"`
}

type CalculateResponse struct {
	Hours     float64 `json:"hours"`
	Days      float64 `json:"days"`
	Formatted string  `json:"formatted"`
	ItemName  string  `json:"itemName"`
	Price     float64 `json:"price"`
	Comment   string  `json:"comment"`
}

type CalculateWithKarmaResponse struct {
	Hours        float64  `json:"hours"`
	Days         float64  `json:"days"`
	Formatted    string   `json:"formatted"`
	ItemName     string   `json:"itemName"`
	Price        float64  `json:"price"`
	Comment      string   `json:"comment"`
	KarmaScore   float64  `json:"karmaScore"`
	KarmaLevel   string   `json:"karmaLevel"`
	KarmaDelta   float64  `json:"karmaDelta"`
	KarmaReasons []string `json:"karmaReasons"`
}
