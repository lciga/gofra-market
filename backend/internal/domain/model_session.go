package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Модель сессии
type Session struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // Уникальный идентификатор сессии
	UserID    primitive.ObjectID `bson:"user_id"`       // Идентификатор пользователя
	SID       string             `bson:"sid"`           // Идентификатор сессии
	ExpiredAt time.Time          `bson:"expires_at"`    // Время истечения сессии
}
