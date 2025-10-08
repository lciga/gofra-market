// Набор моделей
package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Модель гофера
type Gofer struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // Уникальный идентификатор
	OwnerID   primitive.ObjectID `bson:"owner_id"`      // Идентификатор владельца
	Name      string             `bson:"name"`          // Имя гофера
	Rarity    int                `bson:"rarity"`        // Редкость гофера
	CreatedAt time.Time          `bson:"created_at"`    // Время создания
}
