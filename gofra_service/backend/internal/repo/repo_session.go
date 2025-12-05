package repo

import (
	"Gofra_Market/internal/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Структура для репозитория сессий
type SessionRepo struct {
	c *mongo.Collection // Коллекция
}

// Создание нового репозитория
func NewSessionRepo(c *mongo.Collection) *SessionRepo {
	return &SessionRepo{c: c}
}

// Метод для создания сессии
func (r *SessionRepo) Create(ctx context.Context, s *domain.Session) error {
	_, err := r.c.InsertOne(ctx, s)
	return err
}

// Метод для получения сессии по идентификатору
func (r *SessionRepo) ByID(ctx context.Context, id primitive.ObjectID) (*domain.Session, error) {
	var s domain.Session
	err := r.c.FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Метод для получения сессии по SID
func (r *SessionRepo) BySID(ctx context.Context, sid string) (*domain.Session, error) {
	var s domain.Session
	err := r.c.FindOne(ctx, bson.M{"sid": sid}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Метод для удаления сессий
func (r *SessionRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.c.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Количество активных сессий на текущий момент
func (r *SessionRepo) CountActive(ctx context.Context, now time.Time) (int64, error) {
	filter := bson.M{
		"expires_at": bson.M{"$gt": now},
	}
	return r.c.CountDocuments(ctx, filter)
}
