package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Структура конфига приложения
type Config struct {
	MongoURI       string   // URI MongoDB
	DBName         string   // Имя БД
	MongoUser      string   // Пользователь MongoDB
	MongoPassword  string   // Пароль MongoDB
	LogLevel       string   // Уровень логирования
	ServerPort     int      // Порт сервера
	GinMode        string   // Режим работы Gin
	AllowedOrigins []string // Разрешённые источники CORS
	EditorEmail    string   // Почта редактора сайта
}

// Загрузка конфига приложения
func Load() *Config {
	_ = godotenv.Load()
	cfg := &Config{}

	cfg.MongoURI = os.Getenv("MONGO_URL")
	cfg.MongoUser = os.Getenv("MONGO_USER")
	cfg.MongoPassword = os.Getenv("MONGO_PASSWORD")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.LogLevel = os.Getenv("LOG_LEVEL")
	cfg.EditorEmail = os.Getenv("EDITOR_EMAIL")
	cfg.ServerPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
	cfg.GinMode = os.Getenv("GIN_MODE")
	if cfg.GinMode == "" {
		cfg.GinMode = gin.DebugMode
	}
	if cfg.EditorEmail == "" {
		cfg.EditorEmail = "editor@gofra-market.local"
	}

	if corsEnv := os.Getenv("CORS_ALLOWED_ORIGINS"); corsEnv != "" {
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
