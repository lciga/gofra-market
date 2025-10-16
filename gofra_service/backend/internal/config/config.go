// Пакет содержит конфигурацию приложения. загружаемую из переменных окружения
package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Структура для хранения конфигурации приложения
type Config struct {
	MongoURI       string
	DBName         string
	MongoUser      string
	MongoPassword  string
	LogLevel       string
	ServerPort     int
	GinMode        string
	AllowedOrigins []string
}

// Load загружает конфигурацию из переменных окружения.
// Возвращает указатель на Config и ошибку.
func Load() *Config {
	_ = godotenv.Load()
	cfg := &Config{}

	cfg.MongoURI = os.Getenv("MONGO_URL")
	cfg.MongoUser = os.Getenv("MONGO_USER")
	cfg.MongoPassword = os.Getenv("MONGO_PASSWORD")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.LogLevel = os.Getenv("LOG_LEVEL")
	cfg.ServerPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
	cfg.GinMode = os.Getenv("GIN_MODE")

	if corsEnv := os.Getenv("CORS_ALLOWED_ORIGINS"); corsEnv != "" {
		// split by comma and trim whitespace
		parts := strings.Split(corsEnv, ",")
		cfg.AllowedOrigins = make([]string, 0, len(parts))
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				cfg.AllowedOrigins = append(cfg.AllowedOrigins, trimmed)
			}
		}
	}
	if len(cfg.AllowedOrigins) == 0 {
		cfg.AllowedOrigins = []string{"*"}
	}

	return cfg
}
