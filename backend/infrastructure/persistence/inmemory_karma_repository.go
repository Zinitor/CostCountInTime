package persistence

import (
	"sync"

	"costcountintimeapp/domain/entity"
	"costcountintimeapp/domain/repository"
)

// InMemoryKarmaRepository реализует KarmaRepository с использованием in-memory хранилища
type InMemoryKarmaRepository struct {
	events map[string][]*entity.KarmaEvent
	mu     sync.RWMutex
	nextID int64
}

// NewInMemoryKarmaRepository создаёт новый InMemoryKarmaRepository
func NewInMemoryKarmaRepository() *InMemoryKarmaRepository {
	return &InMemoryKarmaRepository{
		events: make(map[string][]*entity.KarmaEvent),
		nextID: 1,
	}
}

// Save сохраняет событие кармы
func (r *InMemoryKarmaRepository) Save(event *entity.KarmaEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if event.ID == 0 {
		event.ID = r.nextID
		r.nextID++
	}

	r.events[event.UserID] = append(r.events[event.UserID], event)
	return nil
}

// GetLast возвращает последнее событие кармы для пользователя
func (r *InMemoryKarmaRepository) GetLast(userID string) (*entity.KarmaEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := r.events[userID]
	if len(events) == 0 {
		return nil, nil
	}

	// Находим последнее событие по timestamp
	var last *entity.KarmaEvent
	for _, event := range events {
		if last == nil || event.Timestamp > last.Timestamp {
			last = event
		}
	}

	return last, nil
}

// GetAllByUser возвращает все события кармы для пользователя
func (r *InMemoryKarmaRepository) GetAllByUser(userID string) ([]*entity.KarmaEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := r.events[userID]
	result := make([]*entity.KarmaEvent, len(events))
	copy(result, events)
	return result, nil
}

// DeleteByItemAndTimestamp удаляет событие по имени предмета и timestamp
func (r *InMemoryKarmaRepository) DeleteByItemAndTimestamp(userID, itemName string, price float64, timestamp int64) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	events := r.events[userID]
	var newEvents []*entity.KarmaEvent
	var deleted int64

	for _, event := range events {
		if event.ItemName == itemName &&
			event.Price == price &&
			abs(event.Timestamp-timestamp) < 100 {
			deleted++
			continue
		}
		newEvents = append(newEvents, event)
	}

	r.events[userID] = newEvents
	return deleted, nil
}

// DeleteAllByUser удаляет все события кармы пользователя
func (r *InMemoryKarmaRepository) DeleteAllByUser(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events[userID] = make([]*entity.KarmaEvent, 0)
	return nil
}

// GetTotalScore возвращает общий счёт кармы пользователя
func (r *InMemoryKarmaRepository) GetTotalScore(userID string) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := r.events[userID]
	var score float64
	for _, event := range events {
		score += event.ScoreDelta
	}
	return score, nil
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// Ensure InMemoryKarmaRepository implements repository.KarmaRepository
var _ repository.KarmaRepository = (*InMemoryKarmaRepository)(nil)

// InitializeDefaultUser инициализирует пользователя по умолчанию
func (r *InMemoryKarmaRepository) InitializeDefaultUser() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.events["default"] == nil {
		r.events["default"] = make([]*entity.KarmaEvent, 0)
	}
}

// GetEventsCount возвращает количество событий для тестирования
func (r *InMemoryKarmaRepository) GetEventsCount(userID string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.events[userID])
}

// Clear очищает все данные для тестирования
func (r *InMemoryKarmaRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = make(map[string][]*entity.KarmaEvent)
	r.nextID = 1
}

// SetTimeFunc устанавливает функцию для получения текущего времени (для тестов)
func (r *InMemoryKarmaRepository) SetTimeFunc(timeFunc func() int64) {
	// Может быть использовано для тестов
}

// GetEventsForUser возвращает события пользователя для тестирования
func (r *InMemoryKarmaRepository) GetEventsForUser(userID string) []*entity.KarmaEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	events := r.events[userID]
	result := make([]*entity.KarmaEvent, len(events))
	copy(result, events)
	return result
}

// AddEventWithTimestamp добавляет событие с указанным timestamp (для тестов)
func (r *InMemoryKarmaRepository) AddEventWithTimestamp(userID, itemName string, price, hoursNeeded, scoreDelta float64, reason string, timestamp int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	event := &entity.KarmaEvent{
		ID:          r.nextID,
		UserID:      userID,
		ItemName:    itemName,
		Price:       price,
		HoursNeeded: hoursNeeded,
		ScoreDelta:  scoreDelta,
		Reason:      reason,
		Timestamp:   timestamp,
	}
	r.nextID++

	r.events[userID] = append(r.events[userID], event)
	return nil
}
