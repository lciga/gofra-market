package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Метаданные изображения
type ImageMeta struct {
	Kind         string     `bson:"kind"`                        // Тип загрузки (например, "url", "upload")
	SourceURL    *string    `bson:"source_url,omitempty"`        // Исходный URL (если применимо)
	FetchedAt    *time.Time `bson:"fetched_at,omitempty"`        // Время получения изображения (если применимо)
	ContentType  *string    `bson:"content_type,omitempty"`      // MIME-тип изображения (если применимо)
	DebugSnippet *string    `bson:"debug_snippet_b64,omitempty"` // Отладочный фрагмент (если применимо)
	ImageData    *string    `bson:"image_data,omitempty"`        // Полное изображение в base64 (для файловых загрузок)
}

// Модель листинга
type Listing struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`      // Уникальный идентификатор
	GoferID     primitive.ObjectID `bson:"gofer_id"`           // Идентификатор гофера
	SellerID    primitive.ObjectID `bson:"seller_id"`          // Идентификатор продавца
	Price       int64              `bson:"price"`              // Цена в минимальных единицах валюты (например, центы)
	IsSold      bool               `bson:"is_sold"`            // Статус продажи
	BuyerID     primitive.ObjectID `bson:"buyer_id,omitempty"` // Идентификатор покупателя (если продано)
	Description string             `bson:"description"`        // Описание листинга чекер вкладывает флаг
	Image       ImageMeta          `bson:"image"`              // Метаданные изображения
	CreatedAt   time.Time          `bson:"created_at"`         // Время создания
}
