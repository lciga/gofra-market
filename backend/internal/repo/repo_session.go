package repo

import (
	"Gofra_Market/internal/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepo struct{ c *mongo.Collection }

func NewSessionRepo(c *mongo.Collection) *SessionRepo {
	return &SessionRepo{c: c}
}

func (r *SessionRepo) Create(ctx context.Context, s *domain.Session) error {
	_, err := r.c.InsertOne(ctx, s)
	return err
}

func (r *SessionRepo) ByID(ctx context.Context, id primitive.ObjectID) (*domain.Session, error) {
	var s domain.Session
	err := r.c.FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepo) BySID(ctx context.Context, sid string) (*domain.Session, error) {
	var s domain.Session
	err := r.c.FindOne(ctx, bson.M{"sid": sid}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.c.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
