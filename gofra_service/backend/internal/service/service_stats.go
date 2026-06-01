package service

import (
	"context"
	"time"
)

// Интерфейс репозитория сессий для статистики
type StatisticsSessionRepo interface {
	CountActive(ctx context.Context, now time.Time) (int64, error)
	CountTotal(ctx context.Context) (int64, error)
}

// Структура статистики посещений
type VisitStatistics struct {
	ActiveUsers int64 `json:"active_users"` // Количество активных пользователей
	TotalVisits int64 `json:"total_visits"` // Количество сохранённых входов
}

// Структура сервиса статистики
type StatisticsService struct {
	sessions StatisticsSessionRepo // Репозиторий сессий
}

// Создание нового сервиса статистики
func NewStatisticsService(repo StatisticsSessionRepo) *StatisticsService {
	return &StatisticsService{sessions: repo}
}

// Метод получения количества активных пользователей
func (s *StatisticsService) ActiveUsers(ctx context.Context) (int64, error) {
	return s.sessions.CountActive(ctx, time.Now())
}

// Метод получения статистики посещений
func (s *StatisticsService) Visits(ctx context.Context) (*VisitStatistics, error) {
	activeUsers, err := s.sessions.CountActive(ctx, time.Now())
	if err != nil {
		return nil, err
	}

	totalVisits, err := s.sessions.CountTotal(ctx)
	if err != nil {
		return nil, err
	}

	return &VisitStatistics{
		ActiveUsers: activeUsers,
		TotalVisits: totalVisits,
	}, nil
}
