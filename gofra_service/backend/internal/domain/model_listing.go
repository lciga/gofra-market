package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageMeta struct {
	Kind         string     `bson:"kind"`                        
	SourceURL    *string    `bson:"source_url,omitempty"`        
	FetchedAt    *time.Time `bson:"fetched_at,omitempty"`        
	ContentType  *string    `bson:"content_type,omitempty"`      
	DebugSnippet *string    `bson:"debug_snippet_b64,omitempty"` 
	ImageData    *string    `bson:"image_data,omitempty"`        

type Listing struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`      
	GoferID     primitive.ObjectID `bson:"gofer_id"`           
	SellerID    primitive.ObjectID `bson:"seller_id"`         
	Price       int64              `bson:"price"`              
	IsSold      bool               `bson:"is_sold"`            
	BuyerID     primitive.ObjectID `bson:"buyer_id,omitempty"` 
	Description string             `bson:"description"`        
	Image       ImageMeta          `bson:"image"`             
	CreatedAt   time.Time          `bson:"created_at"`         
}
