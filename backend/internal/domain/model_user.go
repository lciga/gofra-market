package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Модель пользователя
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // Уникальный идентификатор
	Login     string             `bson:"login"`         // Логин пользователя
	PassHash  []byte             `bson:"pass_hash"`     // Хэш пароля
	Balance   int64              `bson:"balance"`       // Баланс пользователя в минимальных единицах валюты (например, центы)
	CreatedAt time.Time          `bson:"created_at"`    // Время создания
}
