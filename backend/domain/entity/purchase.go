package entity

// Purchase представляет покупку
type Purchase struct {
	ID          int64
	UserID      string
	ItemName    string
	Price       float64
	HoursNeeded float64
	Timestamp   int64
	TagID       string // ID выбранного тега (может быть пустым)
}

// NewPurchase создаёт новую покупку
func NewPurchase(userID, itemName string, price, hoursNeeded float64, timestamp int64) *Purchase {
	return &Purchase{
		UserID:      userID,
		ItemName:    itemName,
		Price:       price,
		HoursNeeded: hoursNeeded,
		Timestamp:   timestamp,
		TagID:       "",
	}
}

// NewPurchaseWithTag создаёт новую покупку с тегом
func NewPurchaseWithTag(userID, itemName string, price, hoursNeeded float64, timestamp int64, tagID string) *Purchase {
	return &Purchase{
		UserID:      userID,
		ItemName:    itemName,
		Price:       price,
		HoursNeeded: hoursNeeded,
		Timestamp:   timestamp,
		TagID:       tagID,
	}
}
