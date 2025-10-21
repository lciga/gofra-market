package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gofer struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` 
	OwnerID   primitive.ObjectID `bson:"owner_id"`      
	Name      string             `bson:"name"`          
	Rarity    int                `bson:"rarity"`        
	CreatedAt time.Time          `bson:"created_at"`   
}
