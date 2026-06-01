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

const (
	RoleAdmin  = "admin"
	RoleEditor = "editor"
	RoleSystem = "system"
)

// Структура сервиса аутентификации
type AuthService struct {
	users    *repo.UserRepo    // Репозиторий пользователей
	sessions *repo.SessionRepo // Репозиторий сессий
	cookie   string            // Имя cookie сессии
}

// Создание нового сервиса аутентификации
func NewAuthService(u *repo.UserRepo, s *repo.SessionRepo, c string) *AuthService {
	return &AuthService{users: u, sessions: s, cookie: c}
}

// Генерация SID сессии
func generateSID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Нормализация роли пользователя
func normalizeRole(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case RoleAdmin:
		return RoleAdmin
	case RoleSystem:
		return RoleSystem
	default:
		return RoleEditor
	}
}

// Метод регистрации пользователя
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
		Role:      RoleEditor,
		Balance:   int64(100),
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

// Метод входа пользователя
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
	user.Role = normalizeRole(user.Role)

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

// Метод получения текущего пользователя
func (a *AuthService) Me(ctx context.Context, userID primitive.ObjectID) (*domain.User, error) {
	user, err := a.users.ByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.Role = normalizeRole(user.Role)
	return user, nil
}
