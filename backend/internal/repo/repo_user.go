package repo

import (
	"Gofra_Market/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct{ c *mongo.Collection }

func (r *UserRepo) Create(ctx context.Context, u *domain.User) (primitive.ObjectID, error) {
	res, err := r.c.InsertOne(ctx, u)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}
