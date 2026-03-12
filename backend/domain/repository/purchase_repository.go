package repository

import "costcountintimeapp/domain/entity"

// PurchaseRepository определяет интерфейс для работы с покупками
type PurchaseRepository interface {
	// Save сохраняет покупку
	Save(purchase *entity.Purchase) error

	// GetAllByUser возвращает все покупки пользователя
	GetAllByUser(userID string) ([]*entity.Purchase, error)

	// DeleteByItemAndTimestamp удаляет покупку по имени и timestamp
	DeleteByItemAndTimestamp(userID, itemName string, price float64, timestamp int64) (int64, error)
}
