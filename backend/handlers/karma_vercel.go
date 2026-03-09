//go:build vercel
// +build vercel

package handlers

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	karmaDBOnce sync.Once
	karmaMemory = make(map[string][]KarmaEvent)
	karmaMutex  sync.RWMutex
)

// InitKarmaDB for Vercel - uses in-memory storage
func InitKarmaDB(dbPath string) error {
	karmaDBOnce.Do(func() {
		// Initialize with empty memory store
		karmaMemory["default"] = []KarmaEvent{}
	})
	return nil
}

// In-memory implementations for Vercel
func saveKarmaEventVercel(event *KarmaEvent) error {
	karmaMutex.Lock()
	defer karmaMutex.Unlock()
	
	events := karmaMemory[event.UserID]
	event.ID = int64(len(events) + 1)
	karmaMemory[event.UserID] = append(events, *event)
	return nil
}

func getLastKarmaEventVercel() (*KarmaEvent, error) {
	karmaMutex.RLock()
	defer karmaMutex.RUnlock()
	
	events := karmaMemory["default"]
	if len(events) == 0 {
		return nil, nil
	}
	event := events[len(events)-1]
	return &event, nil
}

func getCurrentKarmaScoreVercel() float64 {
	karmaMutex.RLock()
	defer karmaMutex.RUnlock()
	
	events := karmaMemory["default"]
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

func getKarmaStatsVercel() (KarmaStats, error) {
	karmaMutex.RLock()
	defer karmaMutex.RUnlock()
	
	events := karmaMemory["default"]
	stats := KarmaStats{}
	
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

func resetKarmaVercel() error {
	karmaMutex.Lock()
	defer karmaMutex.Unlock()
	karmaMemory["default"] = []KarmaEvent{}
	return nil
}

func deleteKarmaEventVercel(itemName string, price float64, timestamp int64) (int64, error) {
	karmaMutex.Lock()
	defer karmaMutex.Unlock()
	
	events := karmaMemory["default"]
	newEvents := []KarmaEvent{}
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
	
	karmaMemory["default"] = newEvents
	return deleted, nil
}
