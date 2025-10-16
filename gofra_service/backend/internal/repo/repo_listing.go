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
	// Generate ID if not set
	if l.ID.IsZero() {
		l.ID = primitive.NewObjectID()
	}
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

func (r *ListingRepo) UpdateImageMeta(ctx context.Context, id primitive.ObjectID, url *string, ct *string, at *time.Time, b64 *string, imageData *string) error {
	// Update nested image fields to match domain.ImageMeta inside Listing.image
	set := bson.M{}
	if url != nil {
		set["image.source_url"] = *url
	}
	if ct != nil {
		set["image.content_type"] = *ct
	}
	if at != nil {
		set["image.fetched_at"] = *at
	}
	if b64 != nil {
		set["image.debug_snippet_b64"] = *b64
	}
	if imageData != nil {
		set["image.image_data"] = *imageData
	}
	update := bson.M{"$set": set}
	_, err := r.c.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Уязвимая функция: запрос кладётся прямо в агрегацию без валидации
// Используем $lookup для присоединения gofers, чтобы можно было фильтровать по gofer.name
func (r *ListingRepo) FindCards(ctx context.Context, raw map[string]any, limit, skip int64, sort bson.D) (cur *mongo.Cursor, total int64, err error) {
	// Build aggregation pipeline
	pipeline := mongo.Pipeline{}
	
	// Stage 1: Match with user filter BEFORE lookup (УЯЗВИМОСТЬ!)
	// Это позволяет фильтровать по полям коллекции listings (включая description)
	if len(raw) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: raw}})
	}
	
	// Stage 2: Lookup gofers
	pipeline = append(pipeline, bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "gofers"},
		{Key: "localField", Value: "gofer_id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "gofer"},
	}}})
	
	// Stage 3: Unwind gofer array (converts array to single object)
	pipeline = append(pipeline, bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$gofer"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}})

	// Stage 4: Sort
	if len(sort) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: sort}})
	}

	// Stage 5: Pagination
	pipeline = append(pipeline,
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
	)

	// Execute aggregation
	cur, err = r.c.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}

	// Count total matching documents (before pagination)
	countPipeline := mongo.Pipeline{}
	
	// Apply match BEFORE lookup for consistency
	if len(raw) > 0 {
		countPipeline = append(countPipeline, bson.D{{Key: "$match", Value: raw}})
	}
	
	// Same lookup and unwind
	countPipeline = append(countPipeline, bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "gofers"},
		{Key: "localField", Value: "gofer_id"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "gofer"},
	}}})
	countPipeline = append(countPipeline, bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$gofer"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}})
	countPipeline = append(countPipeline, bson.D{{Key: "$count", Value: "total"}})

	countCur, err := r.c.Aggregate(ctx, countPipeline)
	if err != nil {
		return cur, 0, nil // Return cursor but no count
	}
	defer countCur.Close(ctx)

	var countResult []struct {
		Total int64 `bson:"total"`
	}
	if err := countCur.All(ctx, &countResult); err == nil && len(countResult) > 0 {
		total = countResult[0].Total
	}

	return cur, total, nil
}

func (r *ListingRepo) ByUser(ctx context.Context, userID primitive.ObjectID) ([]*domain.Listing, error) {
	// Find all listings where user is seller OR buyer
	filter := bson.M{
		"$or": []bson.M{
			{"seller_id": userID},
			{"buyer_id": userID},
		},
	}
	cur, err := r.c.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var listings []*domain.Listing
	if err := cur.All(ctx, &listings); err != nil {
		return nil, err
	}
	return listings, nil
}
