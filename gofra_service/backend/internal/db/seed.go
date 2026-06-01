package db

import (
	"Gofra_Market/internal/domain"
	"Gofra_Market/internal/logger"
	"Gofra_Market/internal/service"
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Создание данных при инициализации
func SeedInitialData(ctx context.Context, db *mongo.Database) error {
	logger.Info("Checking if initial seed is needed", logrus.Fields{})

	if err := ensureAuthSeed(ctx, db); err != nil {
		return err
	}

	metaColl := db.Collection("_seed_meta")
	var meta bson.M
	err := metaColl.FindOne(ctx, bson.M{"_id": "initial_seed"}).Decode(&meta)
	if err == nil {
		logger.Info("Initial seed already applied, skipping", logrus.Fields{})
		return nil
	}
	if err != mongo.ErrNoDocuments {
		return fmt.Errorf("checking seed meta: %w", err)
	}

	logger.Info("Applying initial seed data", logrus.Fields{})

	usersColl := db.Collection("users")
	gofersColl := db.Collection("gofers")
	listingsColl := db.Collection("listings")

	passHash, err := bcrypt.GenerateFromPassword([]byte("system123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	systemUser := domain.User{
		ID:        primitive.NewObjectID(),
		Login:     "system_seller",
		PassHash:  passHash,
		Role:      service.RoleSystem,
		Balance:   100,
		CreatedAt: time.Now(),
	}

	if _, err := usersColl.UpdateOne(
		ctx,
		bson.M{"login": systemUser.Login},
		bson.M{"$setOnInsert": systemUser},
		options.Update().SetUpsert(true),
	); err != nil {
		return fmt.Errorf("upserting system user: %w", err)
	}

	logger.Info("Created system user", logrus.Fields{"login": systemUser.Login, "id": systemUser.ID.Hex()})

	goferNames := []struct {
		name   string
		rarity int
		price  int64
	}{
		{"Common Gopher", 1, 100},
		{"Rare Gopher", 2, 500},
		{"Epic Gopher", 3, 2000},
		{"Legendary Gopher", 3, 5000},
		{"Mythical Gopher", 3, 100},
	}

	for i, gn := range goferNames {
		gofer := domain.Gofer{
			ID:        primitive.NewObjectID(),
			OwnerID:   systemUser.ID,
			Name:      gn.name,
			Rarity:    gn.rarity,
			CreatedAt: time.Now().Add(time.Duration(-i) * time.Minute),
		}

		if _, err := gofersColl.InsertOne(ctx, gofer); err != nil {
			return fmt.Errorf("inserting gofer %s: %w", gofer.Name, err)
		}

		listing := domain.Listing{
			ID:          primitive.NewObjectID(),
			GoferID:     gofer.ID,
			SellerID:    systemUser.ID,
			Price:       gn.price,
			IsSold:      false,
			BuyerID:     primitive.NilObjectID,
			Description: fmt.Sprintf("Прекрасный гофер %s доступен для покупки! Редкость: %d", gofer.Name, gofer.Rarity),
			Image: domain.ImageMeta{
				Kind: "upload",
			},
			CreatedAt: gofer.CreatedAt,
		}

		if _, err := listingsColl.InsertOne(ctx, listing); err != nil {
			return fmt.Errorf("inserting listing for %s: %w", gofer.Name, err)
		}

		logger.Info("Created gofer and listing", logrus.Fields{
			"gofer":   gofer.Name,
			"rarity":  gofer.Rarity,
			"price":   gn.price,
			"listing": listing.ID.Hex(),
		})
	}

	seedMeta := bson.M{
		"_id":        "initial_seed",
		"applied_at": time.Now(),
		"version":    2,
	}
	if _, err := metaColl.InsertOne(ctx, seedMeta); err != nil {
		return fmt.Errorf("saving seed meta: %w", err)
	}

	logger.Info("Initial seed data applied successfully", logrus.Fields{
		"gofers":   len(goferNames),
		"listings": len(goferNames),
	})

	return nil
}

// Создание служебных аккаунтов и ролей
func ensureAuthSeed(ctx context.Context, db *mongo.Database) error {
	usersColl := db.Collection("users")

	if _, err := usersColl.UpdateMany(
		ctx,
		bson.M{"role": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"role": service.RoleEditor}},
	); err != nil {
		return fmt.Errorf("setting default user roles: %w", err)
	}

	adminHash, err := bcrypt.GenerateFromPassword([]byte("admin123456"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing admin password: %w", err)
	}

	adminUser := domain.User{
		ID:        primitive.NewObjectID(),
		Login:     "admin",
		PassHash:  adminHash,
		Role:      service.RoleAdmin,
		Balance:   1000,
		CreatedAt: time.Now(),
	}

	result, err := usersColl.UpdateOne(
		ctx,
		bson.M{"login": adminUser.Login},
		bson.M{
			"$setOnInsert": bson.M{
				"_id":        adminUser.ID,
				"login":      adminUser.Login,
				"balance":    adminUser.Balance,
				"created_at": adminUser.CreatedAt,
			},
			"$set": bson.M{
				"pass_hash": adminUser.PassHash,
				"role":      service.RoleAdmin,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return fmt.Errorf("upserting admin user: %w", err)
	}

	logger.Info("Auth seed checked", logrus.Fields{
		"admin_login": adminUser.Login,
		"matched":     result.MatchedCount,
		"modified":    result.ModifiedCount,
		"upserted":    result.UpsertedCount,
	})
	return nil
}
