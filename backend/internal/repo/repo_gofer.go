package repo

import (
	"Gofra_Market/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GoferRepo struct{ c *mongo.Collection }

func NewGoferRepo(c *mongo.Collection) *GoferRepo {
	return &GoferRepo{c: c}
}

func (r *GoferRepo) Create(ctx context.Context, g *domain.Gofer) error {
	_, err := r.c.InsertOne(ctx, g)
	return err
}

func (r *GoferRepo) ByID(ctx context.Context, id primitive.ObjectID) (*domain.Gofer, error) {
	var g domain.Gofer
	err := r.c.FindOne(ctx, bson.M{"_id": id}).Decode(&g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GoferRepo) TransferOwner(ctx context.Context, goferID, newOwnerID primitive.ObjectID) error {
	_, err := r.c.UpdateOne(ctx, bson.M{"_id": goferID}, bson.M{"$set": bson.M{"owner_id": newOwnerID}})
	return err
}
