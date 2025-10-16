package service

import (
	"Gofra_Market/internal/domain"
	"Gofra_Market/internal/repo"
	"context"
	"encoding/json"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MarketGofer struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Rarity int    `json:"rarity"`
}

type MarketImage struct {
	SourceURL          *string `json:"source_url,omitempty"`
	ContentType        *string `json:"content_type,omitempty"`
	FetchedAt          *string `json:"fetched_at,omitempty"`
	DebugBase64Snippet *string `json:"debug_snippet_b64,omitempty"`
}

type Card struct {
	ID          string      `json:"id"`
	GoferID     string      `json:"gofer_id"`
	SellerID    string      `json:"seller_id"`
	BuyerID     string      `json:"buyer_id,omitempty"`
	Price       int64       `json:"price"`
	IsSold      bool        `json:"is_sold"`
	Description string      `json:"description,omitempty"`
	CreatedAt   string      `json:"created_at"`
	Gofer       MarketGofer `json:"gofer"`
	Image       MarketImage `json:"image"`
}

type MarketService struct {
	listing *repo.ListingRepo
	gofers  *repo.GoferRepo
}

func NewMarketService(l *repo.ListingRepo, g *repo.GoferRepo) *MarketService {
	return &MarketService{listing: l, gofers: g}
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

	cards = []Card{} // Initialize empty slice

	// iterate cursor
	for cur.Next(ctx) {
		var listing domain.Listing
		if err := cur.Decode(&listing); err != nil {
			return nil, 0, err
		}

		gofer, err := s.gofers.ByID(ctx, listing.GoferID)
		if err != nil {
			// Skip this listing if gofer not found instead of failing entirely
			continue
		}

		var buyerID string
		if listing.BuyerID != primitive.NilObjectID {
			buyerID = listing.BuyerID.Hex()
		}

		var fetchedAt *string
		if listing.Image.FetchedAt != nil {
			s := listing.Image.FetchedAt.Format(time.RFC3339)
			fetchedAt = &s
		}

		// Hide description for public market view (only visible to seller/buyer)
		// NoSQL injection still works because it filters at DB level before this
		card := Card{
			ID:          listing.ID.Hex(),
			GoferID:     listing.GoferID.Hex(),
			SellerID:    listing.SellerID.Hex(),
			BuyerID:     buyerID,
			Price:       listing.Price,
			IsSold:      listing.IsSold,
			Description: "", // Always hidden in market listing
			CreatedAt:   listing.CreatedAt.Format(time.RFC3339),
			Gofer: MarketGofer{
				ID:     gofer.ID.Hex(),
				Name:   gofer.Name,
				Rarity: gofer.Rarity,
			},
		}
		card.Image.SourceURL = listing.Image.SourceURL
		card.Image.ContentType = listing.Image.ContentType
		card.Image.FetchedAt = fetchedAt
		card.Image.DebugBase64Snippet = listing.Image.DebugSnippet

		cards = append(cards, card)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}

	return cards, totalRes, nil
}
