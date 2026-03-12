package http

import (
	"encoding/json"
	"net/http"

	"costcountintimeapp/domain/repository"
)

// TagHandler обрабатывает запросы на управление тегами
type TagHandler struct {
	tagRepo repository.TagRepository
}

// TagDTO представляет тег для передачи на клиент
type TagDTO struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ColorStart string `json:"colorStart"`
	ColorEnd   string `json:"colorEnd"`
}

// NewTagHandler создаёт новый TagHandler
func NewTagHandler(tagRepo repository.TagRepository) *TagHandler {
	return &TagHandler{
		tagRepo: tagRepo,
	}
}

// ServeHTTP обрабатывает HTTP запрос
func (h *TagHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tags, err := h.tagRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to get tags", http.StatusInternalServerError)
		return
	}

	// Конвертируем в DTO
	dto := make([]TagDTO, len(tags))
	for i, tag := range tags {
		dto[i] = TagDTO{
			ID:         tag.ID,
			Name:       tag.Name,
			ColorStart: tag.ColorStart,
			ColorEnd:   tag.ColorEnd,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto)
}
