package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Получение идентификатора пользователя из контекста
func currentUserID(c *gin.Context) (primitive.ObjectID, error) {
	v, ok := c.Get("userID")
	if !ok {
		return primitive.NilObjectID, errors.New("unauthenticated")
	}

	uid, ok := v.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("invalid user id in context")
	}

	return uid, nil
}
