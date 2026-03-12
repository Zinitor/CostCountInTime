package repository

import "costcountintimeapp/domain/entity"

// TagRepository определяет интерфейс для работы с тегами
type TagRepository interface {
	// GetAll возвращает все теги
	GetAll() ([]*entity.Tag, error)
	
	// GetByID возвращает тег по ID
	GetByID(id string) (*entity.Tag, error)
}
