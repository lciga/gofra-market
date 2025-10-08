package repo

import (
	"Gofra_Market/internal/domain"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrLoginTaken   = errors.New("login already taken")
)

type UserRepo struct{ c *mongo.Collection }

func NewUserRepo(c *mongo.Collection) *UserRepo {
	return &UserRepo{c: c}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	_, err := r.c.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return ErrLoginTaken
	}
	return err
}

func (r *UserRepo) ByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	err := r.c.FindOne(ctx, bson.M{"login": login}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (r *UserRepo) ByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var user domain.User
	err := r.c.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (r *UserRepo) UpdateBalance(ctx context.Context, id primitive.ObjectID, newBalance int64) error {
	res, err := r.c.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"balance": newBalance}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}
