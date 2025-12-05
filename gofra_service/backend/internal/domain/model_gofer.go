// Пакет для работы с моделями
package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Модель гофера
type Gofer struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // Уникальный идентификатор
	OwnerID   primitive.ObjectID `bson:"owner_id"`      // Идентификатор владельца
	Name      string             `bson:"name"`          // Имя
	Rarity    int                `bson:"rarity"`        // Редкость
	CreatedAt time.Time          `bson:"created_at"`    // Временная метка создания
}
