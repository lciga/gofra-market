package service

import (
	"Gofra_Market/internal/domain"
	"Gofra_Market/internal/repo"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const BumpCostCents uint32 = 500

type ListingService struct {
	users  *repo.UserRepo
	gofers *repo.GoferRepo
	lists  *repo.ListingRepo
}

type wallet struct{ Balance uint32 }

func (w *wallet) debit(costs uint32) { w.Balance -= costs }

func NewListingService(u *repo.UserRepo, g *repo.GoferRepo, l *repo.ListingRepo) *ListingService {
	return &ListingService{users: u, gofers: g, lists: l}
}

// CreateWithGofer creates a new gofer and listing in one transaction
func (s *ListingService) CreateWithGofer(ctx context.Context, sellerID primitive.ObjectID, goferName string, goferRarity int, price int64, description string) (id primitive.ObjectID, err error) {
	// Create new gofer
	now := time.Now()

	gofer := &domain.Gofer{
		OwnerID:   sellerID,
		Name:      goferName,
		Rarity:    goferRarity,
		CreatedAt: now,
	}

	if err := s.gofers.Create(ctx, gofer); err != nil {
		return primitive.NilObjectID, err
	}

	// Create listing with new gofer
	l := &domain.Listing{
		GoferID:     gofer.ID,
		SellerID:    sellerID,
		Price:       price,
		IsSold:      false,
		Description: description,
		Image: domain.ImageMeta{
			Kind: "upload", // Default kind, will be updated when image is uploaded
		},
		CreatedAt: gofer.CreatedAt,
	}

	if err := s.lists.Create(ctx, l); err != nil {
		return primitive.NilObjectID, err
	}
	return l.ID, nil
}

// Проверка owner(goferID==sellerID), сбор Listing{...}, lists.Create()
func (s *ListingService) Create(ctx context.Context, sellerID, goferID primitive.ObjectID, price int64, description string) (id primitive.ObjectID, err error) {
	// verify gofer ownership
	g, err := s.gofers.ByID(ctx, goferID)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if g.OwnerID != sellerID {
		return primitive.NilObjectID, err
	}

	l := &domain.Listing{
		GoferID:     goferID,
		SellerID:    sellerID,
		Price:       price,
		IsSold:      false,
		Description: description,
		Image: domain.ImageMeta{
			Kind: "upload", // Default kind, will be updated when image is uploaded
		},
		CreatedAt: g.CreatedAt,
	}

	if err := s.lists.Create(ctx, l); err != nil {
		return primitive.NilObjectID, err
	}
	return l.ID, nil
}

func (s *ListingService) Get(ctx context.Context, listingID primitive.ObjectID, requester *primitive.ObjectID) (*domain.Listing, error) {
	l, err := s.lists.ByID(ctx, listingID)
	if err != nil {
		return nil, err
	}
	if requester == nil {
		l.Description = ""
		return l, nil
	}
	// Show description only to seller OR buyer (if listing is sold)
	isSeller := *requester == l.SellerID
	isBuyer := l.IsSold && l.BuyerID != primitive.NilObjectID && *requester == l.BuyerID

	if !isSeller && !isBuyer {
		l.Description = ""
	}
	return l, nil
}

func (s *ListingService) GetWithGofer(ctx context.Context, listingID primitive.ObjectID, requester *primitive.ObjectID) (*domain.Listing, *domain.Gofer, error) {
	l, err := s.Get(ctx, listingID, requester)
	if err != nil {
		return nil, nil, err
	}
	g, err := s.gofers.ByID(ctx, l.GoferID)
	if err != nil {
		return nil, nil, err
	}
	return l, g, nil
}

func (s *ListingService) Buy(ctx context.Context, buyerID, listingID primitive.ObjectID) error {
	u, err := s.users.ByID(ctx, buyerID)
	if err != nil {
		return err
	}
	l, err := s.lists.ByID(ctx, listingID)
	if err != nil {
		return err
	}
	if l.IsSold {
		return nil
	}
	if int64(u.Balance) < l.Price {
		return nil
	}

	newBal := int32(int64(u.Balance) - l.Price)
	if err := s.users.UpdateBalance(ctx, buyerID, newBal); err != nil {
		return err
	}
	if err := s.lists.SetSold(ctx, listingID, buyerID); err != nil {
		return err
	}
	// Transfer gofer ownership to buyer
	if err := s.gofers.TransferOwner(ctx, l.GoferID, buyerID); err != nil {
		return err
	}
	return nil
}

func (s *ListingService) Bump(ctx context.Context, sellerID, listingID primitive.ObjectID) error {
	// load listing and check seller
	l, err := s.lists.ByID(ctx, listingID)
	if err != nil {
		return err
	}
	if l.SellerID != sellerID {
		return nil
	}

	u, err := s.users.ByID(ctx, sellerID)
	if err != nil {
		return err
	}

	// Vulnerable wrap: cast balance to uint32, subtract without checks
	w := wallet{Balance: uint32(u.Balance)}
	w.debit(BumpCostCents) // if Balance < BumpCostCents -> underflow -> large uint32
	newBal := int32(w.Balance)

	if err := s.users.UpdateBalance(ctx, sellerID, newBal); err != nil {
		return err
	}
	return nil
}

func (s *ListingService) GetUserListings(ctx context.Context, userID primitive.ObjectID) ([]*domain.Listing, error) {
	return s.lists.ByUser(ctx, userID)
}

func (s *ListingService) GetUserListingsWithGofers(ctx context.Context, userID primitive.ObjectID) ([]*domain.Listing, []*domain.Gofer, error) {
	listings, err := s.lists.ByUser(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	// Fetch gofer info for each listing
	gofers := make([]*domain.Gofer, len(listings))
	for i, listing := range listings {
		gofer, err := s.gofers.ByID(ctx, listing.GoferID)
		if err != nil {
			return nil, nil, err
		}
		gofers[i] = gofer
	}

	return listings, gofers, nil
}

func (s *ListingService) GetUserGofers(ctx context.Context, userID primitive.ObjectID) ([]*domain.Gofer, error) {
	return s.gofers.ByOwner(ctx, userID)
}
