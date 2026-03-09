package handlers

import (
	"sync"
	"time"

	"costcountintimeapp/models"
)

// Simple in-memory karma storage for Vercel
var (
	karmaStore     = make(map[string][]models.KarmaEvent)
	karmaStoreLock sync.RWMutex
)

// InitKarmaDBSimple initializes in-memory storage
func InitKarmaDBSimple(dbPath string) error {
	karmaStoreLock.Lock()
	defer karmaStoreLock.Unlock()
	karmaStore["default"] = []models.KarmaEvent{}
	return nil
}

func saveKarmaEventSimple(event *models.KarmaEvent) error {
	karmaStoreLock.Lock()
	defer karmaStoreLock.Unlock()
	
	events := karmaStore[event.UserID]
	event.ID = int64(len(events) + 1)
	karmaStore[event.UserID] = append(events, *event)
	return nil
}

func getLastKarmaEventSimple() (*models.KarmaEvent, error) {
	karmaStoreLock.RLock()
	defer karmaStoreLock.RUnlock()
	
	events := karmaStore["default"]
	if len(events) == 0 {
		return nil, nil
	}
	event := events[len(events)-1]
	return &event, nil
}

func getCurrentKarmaScoreSimple() float64 {
	karmaStoreLock.RLock()
	defer karmaStoreLock.RUnlock()
	
	events := karmaStore["default"]
	var score float64 = 0
	for _, event := range events {
		score += event.ScoreDelta
	}
	
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	return score
}

func getKarmaStatsSimple() (models.KarmaStats, error) {
	karmaStoreLock.RLock()
	defer karmaStoreLock.RUnlock()
	
	events := karmaStore["default"]
	stats := models.KarmaStats{}
	
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	
	for _, event := range events {
		stats.TotalSpent += event.Price
		stats.TotalHours += event.HoursNeeded
		
		if event.Timestamp >= todayStart {
			stats.PurchasesToday++
			stats.SpentToday += event.Price
		}
	}
	
	if len(events) > 0 {
		stats.AvgPurchase = stats.TotalSpent / float64(len(events))
	}
	
	return stats, nil
}

func resetKarmaSimple() error {
	karmaStoreLock.Lock()
	defer karmaStoreLock.Unlock()
	karmaStore["default"] = []models.KarmaEvent{}
	return nil
}

func deleteKarmaEventSimple(itemName string, price float64, timestamp int64) (int64, error) {
	karmaStoreLock.Lock()
	defer karmaStoreLock.Unlock()
	
	events := karmaStore["default"]
	newEvents := []models.KarmaEvent{}
	var deleted int64 = 0
	
	for _, event := range events {
		if event.ItemName == itemName && 
		   event.Price == price && 
		   event.Timestamp == timestamp {
			deleted++
			continue
		}
		newEvents = append(newEvents, event)
	}
	
	karmaStore["default"] = newEvents
	return deleted, nil
}
