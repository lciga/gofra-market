package service

import (
	"Gofra_Market/internal/domain"
	"Gofra_Market/internal/repo"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const BumpCost uint32 = 500 // Константа для изменения баланса

// Структура сервиса для работы с листингами
type ListingService struct {
	users  *repo.UserRepo    // Репозиторий пользователей
	gofers *repo.GoferRepo   // Репозиторий гоферов
	lists  *repo.ListingRepo // Репозиторий листингов
}

// Структура баланса пользователя
type wallet struct{ Balance uint32 }

// Функция уменьшения баланса
func (w *wallet) debit(costs uint32) { w.Balance -= costs }

// Создание нового сервиса листинга
func NewListingService(u *repo.UserRepo, g *repo.GoferRepo, l *repo.ListingRepo) *ListingService {
	return &ListingService{users: u, gofers: g, lists: l}
}

// Метод для создания листинга
func (s *ListingService) CreateWithGofer(ctx context.Context, sellerID primitive.ObjectID, goferName string, goferRarity int, price int64, description string) (id primitive.ObjectID, err error) {
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

	l := &domain.Listing{
		GoferID:     gofer.ID,
		SellerID:    sellerID,
		Price:       price,
		IsSold:      false,
		Description: description,
		Image: domain.ImageMeta{
			Kind: "upload",
		},
		CreatedAt: gofer.CreatedAt,
	}

	if err := s.lists.Create(ctx, l); err != nil {
		return primitive.NilObjectID, err
	}
	return l.ID, nil
}

// Метод для получение листинга
func (s *ListingService) Get(ctx context.Context, listingID primitive.ObjectID, requester *primitive.ObjectID) (*domain.Listing, error) {
	l, err := s.lists.ByID(ctx, listingID)
	if err != nil {
		return nil, err
	}
	if requester == nil {
		l.Description = ""
		return l, nil
	}
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

// Метод для покупки гофера
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

	newBuyerBal := u.Balance - l.Price
	if err := s.users.UpdateBalance(ctx, buyerID, newBuyerBal); err != nil {
		return err
	}

	seller, err := s.users.ByID(ctx, l.SellerID)
	if err != nil {
		return err
	}
	newSellerBal := seller.Balance + l.Price
	if err := s.users.UpdateBalance(ctx, l.SellerID, newSellerBal); err != nil {
		return err
	}

	if err := s.lists.SetSold(ctx, listingID, buyerID); err != nil {
		return err
	}
	if err := s.gofers.TransferOwner(ctx, l.GoferID, buyerID); err != nil {
		return err
	}
	return nil
}

// Метод для изменения баланса (uint underflow)
func (s *ListingService) Bump(ctx context.Context, sellerID, listingID primitive.ObjectID) error {
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

	// Уязвим для uint underfow при приведении типов
	w := wallet{Balance: uint32(u.Balance)}
	w.debit(BumpCost)
	newBal := int64(w.Balance)

	if err := s.users.UpdateBalance(ctx, sellerID, newBal); err != nil {
		return err
	}
	return nil
}

// Метод для получения листингов пользователя
func (s *ListingService) GetUserListings(ctx context.Context, userID primitive.ObjectID) ([]*domain.Listing, error) {
	return s.lists.ByUser(ctx, userID)
}

// Метод для получения листингов с гоферами
func (s *ListingService) GetUserListingsWithGofers(ctx context.Context, userID primitive.ObjectID) ([]*domain.Listing, []*domain.Gofer, error) {
	listings, err := s.lists.ByUser(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

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

// Метод для получения гоферов пользователя
func (s *ListingService) GetUserGofers(ctx context.Context, userID primitive.ObjectID) ([]*domain.Gofer, error) {
	return s.gofers.ByOwner(ctx, userID)
}
