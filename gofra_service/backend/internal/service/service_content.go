package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"Gofra_Market/internal/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Структура результата подготовки письма редактору
type ContentReview struct {
	EditorEmail string `json:"editor_email"` // Почта редактора
	Subject     string `json:"subject"`      // Тема письма
	Body        string `json:"body"`         // Текст письма
	Mailto      string `json:"mailto"`       // Ссылка mailto
}

// Структура сервиса отправки материалов
type ContentService struct {
	users       *repo.UserRepo // Репозиторий пользователей
	editorEmail string         // Почта редактора
}

// Создание нового сервиса отправки материалов
func NewContentService(users *repo.UserRepo, editorEmail string) *ContentService {
	return &ContentService{users: users, editorEmail: editorEmail}
}

// Метод подготовки письма редактору
func (s *ContentService) PrepareReview(ctx context.Context, userID primitive.ObjectID, title, text string) (*ContentReview, error) {
	user, err := s.users.ByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if normalizeRole(user.Role) != RoleEditor && normalizeRole(user.Role) != RoleAdmin {
		return nil, errors.New("forbidden")
	}

	title = strings.TrimSpace(title)
	text = strings.TrimSpace(text)
	if title == "" || text == "" {
		return nil, errors.New("empty title or text")
	}

	subject := fmt.Sprintf("Материал на проверку: %s", title)
	body := fmt.Sprintf("Автор: %s\n\nЗаголовок:\n%s\n\nТекст материала:\n%s", user.Login, title, text)
	mailto := fmt.Sprintf("mailto:%s?subject=%s&body=%s", s.editorEmail, url.QueryEscape(subject), url.QueryEscape(body))

	return &ContentReview{
		EditorEmail: s.editorEmail,
		Subject:     subject,
		Body:        body,
		Mailto:      mailto,
	}, nil
}
