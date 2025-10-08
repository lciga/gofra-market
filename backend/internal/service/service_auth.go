package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"Gofra_Market/internal/domain"
	"Gofra_Market/internal/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users    *repo.UserRepo
	sessions *repo.SessionRepo
	cookie   string
}

func NewAuthService(u *repo.UserRepo, s *repo.SessionRepo, c string) *AuthService {
	return &AuthService{users: u, sessions: s, cookie: c}
}

// generateSID returns a securely generated random hex string of length 64 (32 bytes).
func generateSID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (a *AuthService) Register(ctx context.Context, login, pass string) (user *domain.User, sid string, err error) {
	login = strings.ToLower(strings.TrimSpace(login))
	if login == "" || pass == "" {
		return nil, "", errors.New("empty login or password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user = &domain.User{
		Login:     login,
		PassHash:  hash,
		Balance:   10000,
		CreatedAt: time.Now(),
	}

	if err := a.users.Create(ctx, user); err != nil {
		return nil, "", err
	}

	sid, err = generateSID()
	if err != nil {
		return nil, "", err
	}

	sess := &domain.Session{
		UserID:    user.ID,
		SID:       sid,
		ExpiredAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := a.sessions.Create(ctx, sess); err != nil {
		return nil, "", err
	}

	return user, sid, nil
}

func (a *AuthService) Login(ctx context.Context, login, pass string) (user *domain.User, sid string, err error) {
	login = strings.ToLower(strings.TrimSpace(login))
	if login == "" || pass == "" {
		return nil, "", errors.New("empty login or password")
	}

	user, err = a.users.ByLogin(ctx, login)
	if err != nil {
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(pass)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	sid, err = generateSID()
	if err != nil {
		return nil, "", err
	}

	sess := &domain.Session{
		UserID:    user.ID,
		SID:       sid,
		ExpiredAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := a.sessions.Create(ctx, sess); err != nil {
		return nil, "", err
	}

	return user, sid, nil
}

func (a *AuthService) Me(ctx context.Context, userID primitive.ObjectID) (*domain.User, error) {
	return a.users.ByID(ctx, userID)
}
