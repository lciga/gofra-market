package repo

import (
	"Gofra_Market/internal/domain"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListingRepo struct{ c *mongo.Collection }

func NewListingRepo(c *mongo.Collection) *ListingRepo {
	return &ListingRepo{c: c}
}

func (r *ListingRepo) Create(ctx context.Context, l *domain.Listing) error {
	_, err := r.c.InsertOne(ctx, l)
	return err
}

func (r *ListingRepo) ByID(ctx context.Context, id primitive.ObjectID) (*domain.Listing, error) {
	var l domain.Listing
	err := r.c.FindOne(ctx, bson.M{"_id": id}).Decode(&l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *ListingRepo) SetSold(ctx context.Context, id, buyer primitive.ObjectID) error {
	// set buyer_id and is_sold flag according to domain.Listing
	update := bson.M{"$set": bson.M{"buyer_id": buyer, "is_sold": true}}
	_, err := r.c.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *ListingRepo) UpdateImageMeta(ctx context.Context, id primitive.ObjectID, url string, ct *string, at *time.Time, b64 *string) error {
	// Update nested image fields to match domain.ImageMeta inside Listing.image
	set := bson.M{"image.source_url": url}
	if ct != nil {
		set["image.content_type"] = *ct
	}
	if at != nil {
		set["image.fetched_at"] = *at
	}
	if b64 != nil {
		set["image.debug_snippet_b64"] = *b64
	}
	update := bson.M{"$set": set}
	_, err := r.c.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Уязвимая функция запрос кладётся прямо в Find()
func (r *ListingRepo) FindCards(ctx context.Context, raw map[string]any, limit, skip int64, sort bson.D) (cur *mongo.Cursor, total int64, err error) {
	cur, err = r.c.Find(ctx, raw, options.Find().SetLimit(limit).SetSkip(skip).SetSort(sort))
	if err != nil {
		return nil, 0, err
	}
	total, err = r.c.CountDocuments(ctx, raw)
	if err != nil {
		return nil, 0, err
	}
	return cur, total, nil
}
