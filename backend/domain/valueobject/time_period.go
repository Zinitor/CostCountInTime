package valueobject

import "errors"

// TimePeriod представляет период времени (час, день, месяц)
type TimePeriod string

const (
	Hour  TimePeriod = "hour"
	Day   TimePeriod = "day"
	Month TimePeriod = "month"
)

// HoursInPeriod возвращает количество часов в периоде
func (p TimePeriod) HoursInPeriod() float64 {
	switch p {
	case Hour:
		return 1
	case Day:
		return 8
	case Month:
		return 168
	default:
		return 0
	}
}

// NewTimePeriod создаёт новый TimePeriod с валидацией
func NewTimePeriod(period string) (TimePeriod, error) {
	switch TimePeriod(period) {
	case Hour, Day, Month:
		return TimePeriod(period), nil
	default:
		return "", errors.New("invalid time period")
	}
}

// IsValid проверяет валидность периода
func (p TimePeriod) IsValid() bool {
	return p == Hour || p == Day || p == Month
}
