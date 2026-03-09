package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"costcountintimeapp/models"

	_ "github.com/mattn/go-sqlite3"
)

var KarmaDB *sql.DB

// Уровни кармы
var karmaLevels = []models.LevelInfo{
	{models.FinancialGod, 0, 20, "Financial God", "😇", "#22c55e"},
	{models.WiseSpender, 20, 40, "Wise Spender", "🧘", "#84cc16"},
	{models.NormalPerson, 40, 60, "Normal Person", "😐", "#eab308"},
	{models.RiskyBuyer, 60, 80, "Risky Buyer", "😬", "#f97316"},
	{models.FinancialDegen, 80, 100, "Financial Degen", "🤡", "#ef4444"},
}

// Инициализация БД
func InitKarmaDB(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS karma_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT DEFAULT 'default',
		item_name TEXT NOT NULL,
		price REAL NOT NULL,
		hours_needed REAL NOT NULL,
		score_delta REAL NOT NULL,
		reason TEXT NOT NULL,
		timestamp INTEGER NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON karma_events(timestamp);
	CREATE INDEX IF NOT EXISTS idx_user ON karma_events(user_id);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		return err
	}

	KarmaDB = db
	return nil
}

// Расчёт дельты кармы
func calculateKarmaDelta(price, hoursNeeded float64, itemName string) (float64, []string) {
	var delta float64
	var reasons []string

	// 1. Цена относительно зарплаты (часы работы)
	if hoursNeeded < 1 {
		delta -= 2
		reasons = append(reasons, "Разумная трата (< 1 часа)")
	} else if hoursNeeded < 8 {
		delta += 0
		reasons = append(reasons, "Нормальная покупка")
	} else if hoursNeeded < 40 {
		delta += 5
		reasons = append(reasons, "Дороговато (8-40 часов)")
	} else if hoursNeeded < 200 {
		delta += 15
		reasons = append(reasons, "Серьёзная трата (40-200 часов)")
	} else {
		delta += 25
		reasons = append(reasons, "Безумие (> 200 часов)")
	}

	// 2. Проверка на "еду" с высокой ценой
	foodKeywords := []string{"салат", "еда", "обед", "ужин", "завтрак", "кофе", "ресторан"}
	for _, keyword := range foodKeywords {
		if containsIgnoreCase(itemName, keyword) && hoursNeeded > 2 {
			delta += 10
			reasons = append(reasons, "Дорогая еда!")
			break
		}
	}

	// 3. Импульсивность (проверяем последнюю покупку)
	lastEvent, err := getLastKarmaEvent()
	if err == nil && lastEvent != nil {
		timeSinceLast := time.Now().Unix() - lastEvent.Timestamp
		if timeSinceLast < 3600 { // < 1 часа
			delta += 10
			reasons = append(reasons, "Импульсивная покупка (< 1 часа)")
		} else if timeSinceLast < 86400 { // < 24 часов
			delta += 3
			reasons = append(reasons, "Покупка в тот же день")
		}
	}

	// 4. Накопленная усталость (сумма за 24 часа)
	stats, _ := getKarmaStats()
	if stats.SpentToday > 50000 { // > 50к за день
		delta += 20
		reasons = append(reasons, "Много потрачено за день")
	}

	// 5. Паттерны (3+ покупки за день)
	if stats.PurchasesToday >= 3 {
		delta += 5
		reasons = append(reasons, "Много покупок сегодня")
	}

	return delta, reasons
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && containsLower(s, substr)
}

func containsLower(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		result[i] = c
	}
	return string(result)
}

// Сохранение события
func saveKarmaEvent(event *models.KarmaEvent) error {
	query := `
	INSERT INTO karma_events (user_id, item_name, price, hours_needed, score_delta, reason, timestamp)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := KarmaDB.Exec(query, event.UserID, event.ItemName, event.Price, event.HoursNeeded, event.ScoreDelta, event.Reason, event.Timestamp)
	return err
}

// Последнее событие
func getLastKarmaEvent() (*models.KarmaEvent, error) {
	query := `
	SELECT id, user_id, item_name, price, hours_needed, score_delta, reason, timestamp
	FROM karma_events
	WHERE user_id = 'default'
	ORDER BY timestamp DESC
	LIMIT 1
	`

	var event models.KarmaEvent
	err := KarmaDB.QueryRow(query).Scan(
		&event.ID, &event.UserID, &event.ItemName, &event.Price,
		&event.HoursNeeded, &event.ScoreDelta, &event.Reason, &event.Timestamp,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// Статистика
func getKarmaStats() (models.KarmaStats, error) {
	var stats models.KarmaStats

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()

	// Всего покупок и сумма
	query := `
	SELECT
		COALESCE(SUM(price), 0),
		COALESCE(SUM(hours_needed), 0),
		COALESCE(AVG(price), 0),
		COUNT(*)
	FROM karma_events
	WHERE user_id = 'default'
	`

	err := KarmaDB.QueryRow(query).Scan(&stats.TotalSpent, &stats.TotalHours, &stats.AvgPurchase, &stats.PurchasesToday)
	if err != nil {
		return stats, err
	}

	// За сегодня
	queryToday := `
	SELECT COUNT(*), COALESCE(SUM(price), 0)
	FROM karma_events
	WHERE user_id = 'default' AND timestamp >= ?
	`

	err = KarmaDB.QueryRow(queryToday, todayStart).Scan(&stats.PurchasesToday, &stats.SpentToday)
	if err != nil {
		return stats, err
	}

	// Импульсивные (в течение часа после предыдущей)
	stats.ImpulseCount = 0

	return stats, nil
}

// Текущий скор
func getCurrentKarmaScore() float64 {
	query := `SELECT COALESCE(SUM(score_delta), 0) FROM karma_events WHERE user_id = 'default'`
	var score float64
	err := KarmaDB.QueryRow(query).Scan(&score)
	if err != nil {
		return 0
	}

	// Decay: уменьшаем на 5% за каждый день
	query = `SELECT MIN(timestamp) FROM karma_events WHERE user_id = 'default'`
	var firstTimestamp int64
	err = KarmaDB.QueryRow(query).Scan(&firstTimestamp)
	if err == nil && firstTimestamp > 0 {
		daysPassed := float64(time.Now().Unix()-firstTimestamp) / 86400
		decay := math.Pow(0.95, daysPassed)
		score *= decay
	}

	// Ограничиваем 0-100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// Определение уровня
func getKarmaLevel(score float64) models.LevelInfo {
	for _, level := range karmaLevels {
		if score >= level.MinScore && score < level.MaxScore {
			return level
		}
	}
	return karmaLevels[len(karmaLevels)-1]
}

// Комментарий на основе уровня
func getKarmaComment(level models.KarmaLevel, delta float64, reasons []string) string {
	comments := map[models.KarmaLevel][]string{
		models.FinancialGod: {
			"😇 Отличный выбор, разумная трата!",
			"🧘 Ты мастер финансов!",
			"💰 Так держать!",
		},
		models.WiseSpender: {
			"👍 Нормально, но можно ещё лучше",
			"🤔 Подумай, можно ли дешевле",
			"📊 В целом разумно",
		},
		models.NormalPerson: {
			"😐 Бывало и лучше",
			"🤷 Нормально для обычного человека",
			"⚖️ Баланс между разумом и эмоциями",
		},
		models.RiskyBuyer: {
			"😬 Чувак, ты там не перегибаешь?",
			"🚩 Красный флаг! Остановись!",
			"💸 К тебе движутся финансовые проблемы",
		},
		models.FinancialDegen: {
			"🤡 Я вызываю финансового консультанта, у тебя проблемы!",
			"🚨 ТРЕВОГА! Твой кошелёк в опасности!",
			"💀💀💀 Это не покупки, это финансовое самоубийство!",
		},
	}

	levelComments := comments[level]
	if len(levelComments) == 0 {
		return ""
	}

	// Добавляем контекст если была большая дельта
	if delta > 15 {
		reason := ""
		if len(reasons) > 0 {
			reason = reasons[0]
		}
		return fmt.Sprintf("%s (%s)", levelComments[0], reason)
	}

	return levelComments[time.Now().Unix()%int64(len(levelComments))]
}

// HTTP Handler
func KarmaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.KarmaUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Расчёт дельты
	delta, reasons := calculateKarmaDelta(req.Price, req.HoursNeeded, req.ItemName)

	// Сохранение события
	event := &models.KarmaEvent{
		UserID:      "default",
		ItemName:    req.ItemName,
		Price:       req.Price,
		HoursNeeded: req.HoursNeeded,
		ScoreDelta:  delta,
		Reason:      joinStrings(reasons, "; "),
		Timestamp:   time.Now().Unix(),
	}

	if err := saveKarmaEvent(event); err != nil {
		http.Error(w, "Failed to save karma event", http.StatusInternalServerError)
		return
	}

	// Получаем новый скор
	newScore := getCurrentKarmaScore()
	level := getKarmaLevel(newScore)
	stats, _ := getKarmaStats()
	comment := getKarmaComment(level.Level, delta, reasons)

	response := models.KarmaUpdateResponse{
		Score:      math.Round(newScore*100) / 100,
		Level:      level.Level,
		ScoreDelta: delta,
		Reasons:    reasons,
		Stats:      stats,
		Comment:    comment,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Текущее состояние кармы
func KarmaStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	score := getCurrentKarmaScore()
	level := getKarmaLevel(score)
	stats, _ := getKarmaStats()

	response := map[string]interface{}{
		"score": score,
		"level": level.Level,
		"info":  level,
		"stats": stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// Сброс кармы
func KarmaResetHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Panic in KarmaResetHandler: %v", rec)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if KarmaDB == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	_, err := KarmaDB.Exec("DELETE FROM karma_events WHERE user_id = 'default'")
	if err != nil {
		log.Printf("Error resetting karma: %v", err)
		http.Error(w, "Failed to reset karma", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"status":  "ok",
		"message": "Карма сброшена",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Удалить событие кармы по item_name и timestamp (для синхронизации с удалением покупки)
func KarmaDeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ItemName  string  `json:"itemName"`
		Price     float64 `json:"price"`
		Timestamp int64   `json:"timestamp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Находим и удаляем событие по приблизительному совпадению
	query := `
	DELETE FROM karma_events
	WHERE user_id = 'default'
	  AND item_name = ?
	  AND ABS(timestamp - ?) < 100
	  AND ABS(price - ?) < 0.01
	`

	result, err := KarmaDB.Exec(query, req.ItemName, req.Timestamp, req.Price)
	if err != nil {
		http.Error(w, "Failed to delete karma event", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	response := map[string]interface{}{
		"status":        "ok",
		"rows_affected": rowsAffected,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
