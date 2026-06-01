package repo

import (
	"Gofra_Market/internal/domain"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Переменные ошибок репозитория пользователей
var (
	ErrUserNotFound = errors.New("user not found")
	ErrLoginTaken   = errors.New("login already taken")
)

// Структура репозитория пользователей
type UserRepo struct {
	c *mongo.Collection // Коллекция пользователей
}

// Создание нового репозитория пользователей
func NewUserRepo(c *mongo.Collection) *UserRepo {
	return &UserRepo{c: c}
}

// Метод создания пользователя
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

// Метод поиска пользователя по логину
func (r *UserRepo) ByLogin(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	err := r.c.FindOne(ctx, bson.M{"login": login}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

// Метод поиска пользователя по идентификатору
func (r *UserRepo) ByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var user domain.User
	err := r.c.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

// Метод получения списка пользователей
func (r *UserRepo) List(ctx context.Context) ([]domain.User, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cur, err := r.c.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	users := make([]domain.User, 0)
	for cur.Next(ctx) {
		var user domain.User
		if err := cur.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, cur.Err()
}

// Метод подсчёта пользователей
func (r *UserRepo) Count(ctx context.Context) (int64, error) {
	return r.c.CountDocuments(ctx, bson.M{})
}

// Метод обновления баланса
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
