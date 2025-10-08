package service

import (
	"Gofra_Market/internal/repo"
	"context"
	"encoding/json"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type Card struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	Teaser string `json:"teaser"`
}

type MarketService struct{ listing *repo.ListingRepo }

func NewMarketService(l *repo.ListingRepo) *MarketService {
	return &MarketService{listing: l}
}

func (s *MarketService) SearchRaw(ctx context.Context, filterJSON string, limit, page int, sort string) (cards []Card, total int64, err error) {
	var raw map[string]any
	if filterJSON != "" {
		if err := json.Unmarshal([]byte(filterJSON), &raw); err != nil {
			return nil, 0, err
		}
	} else {
		raw = map[string]any{}
	}

	// sanitize pagination
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}
	skip := int64((page - 1) * limit)
	lim := int64(limit)

	// sort
	var sortDoc bson.D
	switch strings.ToLower(sort) {
	case "price_asc":
		sortDoc = bson.D{{Key: "price", Value: 1}}
	case "price_desc":
		sortDoc = bson.D{{Key: "price", Value: -1}}
	default:
		sortDoc = bson.D{{Key: "created_at", Value: -1}}
	}

	cur, totalRes, err := s.listing.FindCards(ctx, raw, lim, skip, sortDoc)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	// iterate cursor
	for cur.Next(ctx) {
		var item map[string]any
		if err := cur.Decode(&item); err != nil {
			return nil, 0, err
		}

		// map fields defensively but allow missing fields
		id := ""
		if _id, ok := item["_id"]; ok {
			// try to stringify ObjectID or raw value
			id = toString(_id)
		}
		name := toString(item["name"])
		price := int64(0)
		if p, ok := item["price"].(float64); ok {
			price = int64(p)
		}
		desc := toString(item["description"])
		teaser := desc
		if len(teaser) > 64 {
			teaser = teaser[:64]
		}

		cards = append(cards, Card{ID: id, Name: name, Price: price, Teaser: teaser})
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}

	return cards, totalRes, nil
}

// toString attempts to convert various raw BSON/JSON values to string for display
func toString(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		// fallback to json marshal/unmarshal trick
		b, err := json.Marshal(t)
		if err != nil {
			return ""
		}
		s := string(b)
		return s
	}
}
