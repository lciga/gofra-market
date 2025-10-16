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
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//go:embed migration/*
var migrationsFS embed.FS

// Подключение к БД
func Connect(cfg *config.Config) (*mongo.Client, *mongo.Database, error) {
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

	if cfg.MongoUser != "" {
		// Root user (created via MONGO_INITDB_ROOT_USERNAME) is stored in the "admin" database
		clientOptions.SetAuth(options.Credential{
			Username:   cfg.MongoUser,
			Password:   cfg.MongoPassword,
			AuthSource: "admin",
		})
	}

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
func Migrate(ctx context.Context, cfg *config.Config) error {
	logger.Info("Starting MongoDB migration", logrus.Fields{
		"db": cfg.DBName,
	})

	// Источник миграций (embed)
	src, err := iofs.New(migrationsFS, "migration")
	if err != nil {
		logger.Error(fmt.Errorf("iofs.New failed: %w", err), logrus.Fields{
			"step": "iofs.New",
		})
		return fmt.Errorf("iofs.New: %w", err)
	}

	// Для отладки — перечислим файлы в embed, чтобы понять что включено
	if entries, err := migrationsFS.ReadDir("migration"); err != nil {
		logger.Warnf("failed to read embedded migrations dir: %v", err)
	} else {
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		logger.Info("embedded migrations files", logrus.Fields{"files": names})
	}

	// Подключаемся к БД с таймаутом
	connCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.MongoURI)
	if cfg.MongoUser != "" {
		clientOpts.SetAuth(options.Credential{
			Username:   cfg.MongoUser,
			Password:   cfg.MongoPassword,
			AuthSource: "admin",
		})
	}
	client, err := mongo.Connect(connCtx, clientOpts)
	if err != nil {
		logger.Error(fmt.Errorf("mongo.Connect failed: %w", err), logrus.Fields{
			"step": "mongo.Connect",
			"db":   cfg.DBName,
		})
		return fmt.Errorf("mongo.Connect: %w", err)
	}

	// Всегда попытаться корректно закрыть клиент
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			logger.Error(fmt.Errorf("client.Disconnect failed: %w", err), logrus.Fields{"step": "client.Disconnect"})
		} else {
			logger.Debugf("mongo client disconnected for db %s", cfg.DBName)
		}
	}()

	// Ping — убеждаемся, что БД отвечает
	if err := client.Ping(connCtx, readpref.Primary()); err != nil {
		logger.Error(fmt.Errorf("ping primary failed: %w", err), logrus.Fields{"step": "ping"})
		return fmt.Errorf("ping primary: %w", err)
	}

	// Драйвер migrate
	// Transactions are only supported on replica set members or mongos.
	// Disable TransactionMode for standalone servers (common in local docker setups).
	driver, err := mgm.WithInstance(client, &mgm.Config{
		DatabaseName:         cfg.DBName,
		MigrationsCollection: "gofra",
		TransactionMode:      false,
	})
	if err != nil {
		logger.Error(fmt.Errorf("mgm.WithInstance failed: %w", err), logrus.Fields{"step": "mgm.WithInstance"})
		return fmt.Errorf("mgm.WithInstance: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "mongodb", driver)
	if err != nil {
		logger.Error(fmt.Errorf("migrate.NewWithInstance failed: %w", err), logrus.Fields{"step": "NewWithInstance"})
		return fmt.Errorf("migrate.NewWithInstance: %w", err)
	}

	defer func() {
		if _, cerr := m.Close(); cerr != nil {
			logger.Warnf("m.Close() error: %v", cerr)
		}
	}()

	// Прогоняем вверх
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("Mongo migration: up to date")
			return nil
		}
		logger.Error(fmt.Errorf("m.Up failed: %w", err), logrus.Fields{"step": "m.Up"})
		return fmt.Errorf("m.Up: %w", err)
	}

	logger.Info("Mongo migration: successfully applied all migrations")
	return nil
}

// Штатное закрытие соединения
func Close(client *mongo.Client) error {
	if client == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		logger.Error(fmt.Errorf("client.Disconnect failed: %w", err), logrus.Fields{"step": "Close"})
		return err
	}
	logger.Debugf("Mongo client disconnected cleanly")
	return nil
}
