package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` 
	Login     string             `bson:"login"`         
	PassHash  []byte             `bson:"pass_hash"`    
	Balance   int64              `bson:"balance"`     
	CreatedAt time.Time          `bson:"created_at"`    
}
