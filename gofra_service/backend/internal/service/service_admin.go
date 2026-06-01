package service

import (
	"context"
	"errors"
	"time"

	"Gofra_Market/internal/domain"
	"Gofra_Market/internal/repo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Структура пользователя для админской панели
type AdminUser struct {
	ID        string    `json:"id"`         // Идентификатор пользователя
	Login     string    `json:"login"`      // Логин пользователя
	Role      string    `json:"role"`       // Роль пользователя
	Balance   int64     `json:"balance"`    // Баланс пользователя
	CreatedAt time.Time `json:"created_at"` // Время создания пользователя
}

// Структура статистики админской панели
type AdminStats struct {
	UsersTotal  int64 `json:"users_total"`  // Количество пользователей
	ActiveUsers int64 `json:"active_users"` // Количество активных пользователей
	TotalVisits int64 `json:"total_visits"` // Количество сохранённых входов
}

// Структура ответа админской панели
type AdminDashboard struct {
	Users []AdminUser `json:"users"` // Пользователи системы
	Stats AdminStats  `json:"stats"` // Статистика системы
}

// Структура сервиса админской панели
type AdminService struct {
	users    *repo.UserRepo    // Репозиторий пользователей
	sessions *repo.SessionRepo // Репозиторий сессий
}

// Создание нового сервиса админской панели
func NewAdminService(users *repo.UserRepo, sessions *repo.SessionRepo) *AdminService {
	return &AdminService{users: users, sessions: sessions}
}

// Метод получения данных админской панели
func (s *AdminService) Dashboard(ctx context.Context, requesterID primitive.ObjectID) (*AdminDashboard, error) {
	requester, err := s.users.ByID(ctx, requesterID)
	if err != nil {
		return nil, err
	}
	if normalizeRole(requester.Role) != RoleAdmin {
		return nil, errors.New("forbidden")
	}

	users, err := s.users.List(ctx)
	if err != nil {
		return nil, err
	}

	mappedUsers := make([]AdminUser, 0, len(users))
	for _, user := range users {
		mappedUsers = append(mappedUsers, mapAdminUser(user))
	}

	usersTotal, err := s.users.Count(ctx)
	if err != nil {
		return nil, err
	}

	activeUsers, err := s.sessions.CountActive(ctx, time.Now())
	if err != nil {
		return nil, err
	}

	totalVisits, err := s.sessions.CountTotal(ctx)
	if err != nil {
		return nil, err
	}

	return &AdminDashboard{
		Users: mappedUsers,
		Stats: AdminStats{
			UsersTotal:  usersTotal,
			ActiveUsers: activeUsers,
			TotalVisits: totalVisits,
		},
	}, nil
}

// Преобразование пользователя в ответ админской панели
func mapAdminUser(user domain.User) AdminUser {
	return AdminUser{
		ID:        user.ID.Hex(),
		Login:     user.Login,
		Role:      normalizeRole(user.Role),
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt,
	}
}
