package persistence

import (
	"database/sql"
	"fmt"

	"costcountintimeapp/domain/entity"
	"costcountintimeapp/domain/repository"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteKarmaRepository реализует KarmaRepository с использованием SQLite
type SQLiteKarmaRepository struct {
	db *sql.DB
}

// NewSQLiteKarmaRepository создаёт новый SQLiteKarmaRepository
func NewSQLiteKarmaRepository(dbPath string) (*SQLiteKarmaRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Создаём таблицу
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &SQLiteKarmaRepository{db: db}, nil
}

func createTables(db *sql.DB) error {
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

	_, err := db.Exec(createTable)
	return err
}

// Save сохраняет событие кармы
func (r *SQLiteKarmaRepository) Save(event *entity.KarmaEvent) error {
	query := `
	INSERT INTO karma_events (user_id, item_name, price, hours_needed, score_delta, reason, timestamp)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, event.UserID, event.ItemName, event.Price, event.HoursNeeded, event.ScoreDelta, event.Reason, event.Timestamp)
	return err
}

// GetLast возвращает последнее событие кармы для пользователя
func (r *SQLiteKarmaRepository) GetLast(userID string) (*entity.KarmaEvent, error) {
	query := `
	SELECT id, user_id, item_name, price, hours_needed, score_delta, reason, timestamp
	FROM karma_events
	WHERE user_id = ?
	ORDER BY timestamp DESC
	LIMIT 1
	`

	var event entity.KarmaEvent
	err := r.db.QueryRow(query, userID).Scan(
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

// GetAllByUser возвращает все события кармы для пользователя
func (r *SQLiteKarmaRepository) GetAllByUser(userID string) ([]*entity.KarmaEvent, error) {
	query := `
	SELECT id, user_id, item_name, price, hours_needed, score_delta, reason, timestamp
	FROM karma_events
	WHERE user_id = ?
	ORDER BY timestamp DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*entity.KarmaEvent
	for rows.Next() {
		var event entity.KarmaEvent
		if err := rows.Scan(
			&event.ID, &event.UserID, &event.ItemName, &event.Price,
			&event.HoursNeeded, &event.ScoreDelta, &event.Reason, &event.Timestamp,
		); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, rows.Err()
}

// DeleteByItemAndTimestamp удаляет событие по имени предмета и timestamp
func (r *SQLiteKarmaRepository) DeleteByItemAndTimestamp(userID, itemName string, price float64, timestamp int64) (int64, error) {
	query := `
	DELETE FROM karma_events
	WHERE user_id = ?
	  AND item_name = ?
	  AND ABS(timestamp - ?) < 100
	  AND ABS(price - ?) < 0.01
	`

	result, err := r.db.Exec(query, userID, itemName, timestamp, price)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// DeleteAllByUser удаляет все события кармы пользователя
func (r *SQLiteKarmaRepository) DeleteAllByUser(userID string) error {
	_, err := r.db.Exec("DELETE FROM karma_events WHERE user_id = ?", userID)
	return err
}

// GetTotalScore возвращает общий счёт кармы пользователя
func (r *SQLiteKarmaRepository) GetTotalScore(userID string) (float64, error) {
	query := `SELECT COALESCE(SUM(score_delta), 0) FROM karma_events WHERE user_id = ?`
	var score float64
	err := r.db.QueryRow(query, userID).Scan(&score)
	if err != nil {
		return 0, err
	}
	return score, nil
}

// Close закрывает соединение с БД
func (r *SQLiteKarmaRepository) Close() error {
	return r.db.Close()
}

// Ensure SQLiteKarmaRepository implements repository.KarmaRepository
var _ repository.KarmaRepository = (*SQLiteKarmaRepository)(nil)
