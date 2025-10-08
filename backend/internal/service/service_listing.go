package service

import (
	"Gofra_Market/internal/domain"
	"Gofra_Market/internal/repo"
	"context"

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
		CreatedAt:   g.CreatedAt,
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
	if *requester != l.SellerID && (l.BuyerID == primitive.NilObjectID || *requester != l.BuyerID) {
		l.Description = ""
	}
	return l, nil
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
	if u.Balance < l.Price {
		return nil
	}

	if err := s.users.UpdateBalance(ctx, buyerID, u.Balance-l.Price); err != nil {
		return err
	}
	if err := s.lists.SetSold(ctx, listingID, buyerID); err != nil {
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
	u.Balance = int64(w.Balance)

	if err := s.users.UpdateBalance(ctx, sellerID, u.Balance); err != nil {
		return err
	}
	return nil
}
