package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Структура метаданных для изображений
type ImageMeta struct {
	Kind         string     `bson:"kind"`                        // Тип загрузки: upload, url
	SourceURL    *string    `bson:"source_url,omitempty"`        // URL источника
	FetchedAt    *time.Time `bson:"fetched_at,omitempty"`        // Время загрузки
	ContentType  *string    `bson:"content_type,omitempty"`      // Тип содержимого
	DebugSnippet *string    `bson:"debug_snippet_b64,omitempty"` // Сниппет для дебага
	ImageData    *string    `bson:"image_data,omitempty"`        // Данные изображения
}

// Структура листинга
type Listing struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`      // Уникальный идентификатор
	GoferID     primitive.ObjectID `bson:"gofer_id"`           // Идентификатор гофера
	SellerID    primitive.ObjectID `bson:"seller_id"`          // Идентификатор продавца
	Price       int64              `bson:"price"`              // Цена
	IsSold      bool               `bson:"is_sold"`            // Флаг для проверки на факт продажи
	BuyerID     primitive.ObjectID `bson:"buyer_id,omitempty"` // Идентификатор покупателя
	Description string             `bson:"description"`        // Описание (здесь лежит флаг)
	Image       ImageMeta          `bson:"image"`              // Изображение
	CreatedAt   time.Time          `bson:"created_at"`         // Временная метка создания
}
