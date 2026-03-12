package persistence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"costcountintimeapp/domain/entity"
	"costcountintimeapp/domain/repository"
)

// BlobKarmaRepository реализует KarmaRepository с использованием Vercel Blob
type BlobKarmaRepository struct {
	token string
}

// KarmaData представляет все данные кармы для пользователя
type KarmaData struct {
	Events []*entity.KarmaEvent `json:"events"`
}

// NewBlobKarmaRepository создаёт новый BlobKarmaRepository
func NewBlobKarmaRepository() (*BlobKarmaRepository, error) {
	token := os.Getenv("BLOB_READ_WRITE_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("BLOB_READ_WRITE_TOKEN not set")
	}

	return &BlobKarmaRepository{
		token: token,
	}, nil
}

func (r *BlobKarmaRepository) getBlobPath(userID string) string {
	return fmt.Sprintf("karma/%s.json", userID)
}

func (r *BlobKarmaRepository) loadKarmaData(userID string) (*KarmaData, error) {
	ctx := context.Background()
	req, err := r.newBlobRequest(ctx, "GET", r.getBlobPath(userID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch blob: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &KarmaData{Events: make([]*entity.KarmaEvent, 0)}, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("blob API error: %d - %s", resp.StatusCode, string(body))
	}

	var data KarmaData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode karma data: %w", err)
	}

	return &data, nil
}

func (r *BlobKarmaRepository) saveKarmaData(userID string, data *KarmaData) error {
	ctx := context.Background()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal karma data: %w", err)
	}

	req, err := r.newBlobRequest(ctx, "PUT", r.getBlobPath(userID), bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to save blob: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("blob API error: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

func (r *BlobKarmaRepository) newBlobRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("https://api.vercel.com/v1/blob/%s", path)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+r.token)
	return req, nil
}

// Save сохраняет событие кармы
func (r *BlobKarmaRepository) Save(event *entity.KarmaEvent) error {
	data, err := r.loadKarmaData(event.UserID)
	if err != nil {
		return err
	}

	if event.ID == 0 {
		// Generate new ID
		event.ID = time.Now().UnixNano()
	}

	data.Events = append(data.Events, event)
	return r.saveKarmaData(event.UserID, data)
}

// GetLast возвращает последнее событие кармы для пользователя
func (r *BlobKarmaRepository) GetLast(userID string) (*entity.KarmaEvent, error) {
	data, err := r.loadKarmaData(userID)
	if err != nil {
		return nil, err
	}

	if len(data.Events) == 0 {
		return nil, nil
	}

	// Find last event by timestamp
	var last *entity.KarmaEvent
	for _, event := range data.Events {
		if last == nil || event.Timestamp > last.Timestamp {
			last = event
		}
	}

	return last, nil
}

// GetAllByUser возвращает все события кармы для пользователя
func (r *BlobKarmaRepository) GetAllByUser(userID string) ([]*entity.KarmaEvent, error) {
	data, err := r.loadKarmaData(userID)
	if err != nil {
		return nil, err
	}

	return data.Events, nil
}

// DeleteByItemAndTimestamp удаляет событие по имени предмета и timestamp
func (r *BlobKarmaRepository) DeleteByItemAndTimestamp(userID, itemName string, price float64, timestamp int64) (int64, error) {
	data, err := r.loadKarmaData(userID)
	if err != nil {
		return 0, err
	}

	var newEvents []*entity.KarmaEvent
	var deleted int64

	for _, event := range data.Events {
		if event.ItemName == itemName &&
			event.Price == price &&
			absInt64(event.Timestamp-timestamp) < 100 {
			deleted++
			continue
		}
		newEvents = append(newEvents, event)
	}

	data.Events = newEvents
	if err := r.saveKarmaData(userID, data); err != nil {
		return 0, err
	}

	return deleted, nil
}

// DeleteAllByUser удаляет все события кармы пользователя
func (r *BlobKarmaRepository) DeleteAllByUser(userID string) error {
	data := &KarmaData{Events: make([]*entity.KarmaEvent, 0)}
	return r.saveKarmaData(userID, data)
}

// GetTotalScore возвращает общий счёт кармы пользователя
func (r *BlobKarmaRepository) GetTotalScore(userID string) (float64, error) {
	data, err := r.loadKarmaData(userID)
	if err != nil {
		return 0, err
	}

	var score float64
	for _, event := range data.Events {
		score += event.ScoreDelta
	}

	// Apply decay based on days since first event
	if len(data.Events) > 0 {
		var firstTimestamp int64
		for _, event := range data.Events {
			if firstTimestamp == 0 || event.Timestamp < firstTimestamp {
				firstTimestamp = event.Timestamp
			}
		}

		if firstTimestamp > 0 {
			// Could apply decay here if needed
			_ = firstTimestamp // silence unused warning
		}
	}

	return score, nil
}

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// Ensure BlobKarmaRepository implements repository.KarmaRepository
var _ repository.KarmaRepository = (*BlobKarmaRepository)(nil)
