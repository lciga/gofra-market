// Пакет для работы с БДСМ
package db

import (
	"Gofra_Market/internal/config"
	"Gofra_Market/internal/logger"
	"context"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	mgm "github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// go:embed migrations/*.js
var migrationsFS embed.FS

// Подключение к БД
func Connect(cfg config.Config) (*mongo.Client, *mongo.Database, error) {
	if cfg.MongoURI == "" {
		return nil, nil, errors.New("Connect: empty MongoDB URI")
	}
	if cfg.DBName == "" {
		return nil, nil, errors.New("Connect: empty database name")
	}

	// Конфигурация клиента
	clientOptions := options.Client().
		ApplyURI(cfg.MongoURI).
		SetRetryReads(true).
		SetRetryWrites(true).
		SetAppName("gofra-market")

	// Подулючаемся с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("Connect: mongo.Connect: %w", err)
	}

	// Проверяем доступность
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		client.Disconnect(context.Background())
		return nil, nil, fmt.Errorf("Connect: ping primary failed: %w", err)
	}
	db := client.Database(cfg.DBName)
	return client, db, nil
}

// Выполнение миграции
func Migrate(ctx context.Context, cfg config.Config) error {
	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}

	// Клиент для подключения к БД
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))

	// Драйвер migrate
	driver, err := mgm.WithInstance(client, &mgm.Config{
		DatabaseName:         cfg.DBName,
		MigrationsCollection: "gofra",
		TransactionMode:      true,
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", src, "mongodb", driver)
	if err != nil {
		return err
	}

	defer func() {
		m.Close()
	}()

	// Прогоняем вверх
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	logger.Info("Mongo migration: up to date")
	return nil
}

// Штатное закрытие соединения
func Close(client *mongo.Client) error {
	if client == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Disconnect(ctx)
}
