// Пакет содержит конфигурацию приложения. загружаемую из переменных окружения
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Структура для хранения конфигурации приложения
type Config struct {
	MongoURI   string
	DBName     string
	LogLevel   string
	ServerPort int
	GinMode    string
}

// Load загружает конфигурацию из переменных окружения.
// Возвращает указатель на Config и ошибку.
func Load() *Config {
	_ = godotenv.Load()
	cfg := &Config{}

	cfg.MongoURI = os.Getenv("MONGO_URL")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.LogLevel = os.Getenv("LOG_LEVEL")
	cfg.ServerPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
	cfg.GinMode = os.Getenv("GIN_MODE")

	return cfg
}
