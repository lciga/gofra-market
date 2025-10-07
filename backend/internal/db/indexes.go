package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Создание нужных индексов в БД
func EnsureIndexes(db mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Индекс для уникальности users
	if _, err := db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "login", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("ux_users_login"),
	}); err != nil {
		return fmt.Errorf("ensure users.login index: %w", err)
	}

	// Индекс для проверки продажи
	if _, err := db.Collection("listing").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "is_sold", Value: 1}, {Key: "price", Value: 1}},
			Options: options.Index().SetName("ix_listings_is_sold_price"),
		},
		{
			Keys:    bson.D{{Key: "saller_id", Value: 1}},
			Options: options.Index().SetName("ix_listings_seller_id"),
		},
	}); err != nil {
		return fmt.Errorf("ensure listings indexes: %w", err)
	}

	// Проверка TTL сессии
	if _, err := db.Collection("sessions").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0).SetName("ttl_sessions_expires_at"),
	}); err != nil {
		return fmt.Errorf("ensure sessions TTL index: %w", err)
	}

	return nil
}
