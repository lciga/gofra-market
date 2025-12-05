package repo

import (
	"Gofra_Market/internal/domain"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Переменные для ошибок
var (
	ErrUserNotFound = errors.New("user not found")
	ErrLoginTaken   = errors.New("login already taken")
)

// Структура для репозитория пользователей
type UserRepo struct {
	c *mongo.Collection // Коллекция
}

// Создания нового репозитория пользователей
func NewUserRepo(c *mongo.Collection) *UserRepo {
	return &UserRepo{c: c}
}

// Метод для создания пользователя
func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	res, err := r.c.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return ErrLoginTaken
	}
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}
	return nil
}

// Метод для поиска пользователя по логину
func (r *UserRepo) ByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	err := r.c.FindOne(ctx, bson.M{"login": login}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

// Метод для поиска пользователя по идентификатору
func (r *UserRepo) ByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var user domain.User
	err := r.c.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

// Метод для обновления баланса
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
