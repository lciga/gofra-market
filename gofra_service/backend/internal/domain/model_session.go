package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Структура сессии
type Session struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // Уникальный идентификатор
	UserID    primitive.ObjectID `bson:"user_id"`       // Идентификатор пользователя
	SID       string             `bson:"sid"`           // SID сессии
	ExpiredAt time.Time          `bson:"expires_at"`    // Время истечения сессии
}
