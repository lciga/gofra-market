package service

import (
	"Gofra_Market/internal/domain"
	"Gofra_Market/internal/repo"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminService struct{ list *repo.ListingRepo }

func NewAdminService(l *repo.ListingRepo) *AdminService {
	return &AdminService{list: l}
}

func (s *AdminService) ExportListing(ctx context.Context, id primitive.ObjectID) (*domain.Listing, error) {
	return s.list.ByID(ctx, id)
}
