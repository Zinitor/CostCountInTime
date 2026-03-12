package valueobject

import "errors"

// Money представляет денежную сумму
type Money float64

// NewMoney создаёт новый Money с валидацией
func NewMoney(amount float64) (Money, error) {
	if amount <= 0 {
		return 0, errors.New("amount must be positive")
	}
	return Money(amount), nil
}

// ToFloat64 возвращает значение как float64
func (m Money) ToFloat64() float64 {
	return float64(m)
}

// IsPositive проверяет, положительна ли сумма
func (m Money) IsPositive() bool {
	return m > 0
}
