package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Структура пользователя
type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty"` // Уникальный идентификатор

	Login    string `bson:"login"`     // Логин
	PassHash []byte `bson:"pass_hash"` // Хэш пароля
	Role     string `bson:"role"`      // Роль пользователя
	Balance  int64  `bson:"balance"`   // Баланс

	CreatedAt time.Time `bson:"created_at"` // Время создания
}
