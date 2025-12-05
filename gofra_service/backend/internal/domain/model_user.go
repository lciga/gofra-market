package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Структура пользователя
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // Уникальный идентификатор
	Login     string             `bson:"login"`         // Логин
	PassHash  []byte             `bson:"pass_hash"`     // Хэш пароля
	Balance   int64              `bson:"balance"`       // Баланс
	CreatedAt time.Time          `bson:"created_at"`    // Временная метка создания
}
