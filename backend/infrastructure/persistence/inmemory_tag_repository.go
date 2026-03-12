package persistence

import (
	"sync"

	"costcountintimeapp/domain/entity"
	"costcountintimeapp/domain/repository"
)

// InMemoryTagRepository реализует TagRepository с использованием in-memory хранилища
type InMemoryTagRepository struct {
	tags map[string]*entity.Tag
	mu   sync.RWMutex
}

// NewInMemoryTagRepository создаёт новый InMemoryTagRepository
func NewInMemoryTagRepository() *InMemoryTagRepository {
	repo := &InMemoryTagRepository{
		tags: make(map[string]*entity.Tag),
	}
	
	// Инициализируем теги по умолчанию
	for _, tag := range entity.DefaultTags {
		repo.tags[tag.ID] = &tag
	}
	
	return repo
}

// GetAll возвращает все теги
func (r *InMemoryTagRepository) GetAll() ([]*entity.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	result := make([]*entity.Tag, 0, len(r.tags))
	for _, tag := range r.tags {
		result = append(result, tag)
	}
	return result, nil
}

// GetByID возвращает тег по ID
func (r *InMemoryTagRepository) GetByID(id string) (*entity.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	tag, ok := r.tags[id]
	if !ok {
		return nil, nil
	}
	return tag, nil
}

// Ensure InMemoryTagRepository implements repository.TagRepository
var _ repository.TagRepository = (*InMemoryTagRepository)(nil)
