package repository

import (
	"costcountintimeapp/domain/entity"
)

// KarmaRepository определяет интерфейс для работы с событиями кармы
type KarmaRepository interface {
	// Save сохраняет событие кармы
	Save(event *entity.KarmaEvent) error

	// GetLast возвращает последнее событие кармы для пользователя
	GetLast(userID string) (*entity.KarmaEvent, error)

	// GetAllByUser возвращает все события кармы для пользователя
	GetAllByUser(userID string) ([]*entity.KarmaEvent, error)

	// DeleteByItemAndTimestamp удаляет событие по имени предмета и timestamp
	DeleteByItemAndTimestamp(userID, itemName string, price float64, timestamp int64) (int64, error)

	// DeleteAllByUser удаляет все события кармы пользователя
	DeleteAllByUser(userID string) error

	// GetTotalScore возвращает общий счёт кармы пользователя
	GetTotalScore(userID string) (float64, error)
}
