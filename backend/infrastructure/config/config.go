package config

import "os"

// Config содержит конфигурацию приложения
type Config struct {
	Port   string
	DBPath string
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "karma.db"
	}

	return &Config{
		Port:   port,
		DBPath: dbPath,
	}
}
