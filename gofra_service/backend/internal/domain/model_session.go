package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`       
	SID       string             `bson:"sid"`          
	ExpiredAt time.Time          `bson:"expires_at"`   
}
