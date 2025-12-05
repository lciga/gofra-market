package service

import (
	"context"
	"time"
)

// Интерфейс репозитория сессий, необходимый для статистики
//
//go:generate mockgen -destination=../mocks/mock_statistics_repo.go -package=mocks . StatisticsSessionRepo
type StatisticsSessionRepo interface {
	CountActive(ctx context.Context, now time.Time) (int64, error)
}

// Предоставляет агрегированные метрики по пользователям
// для отображения на фронтенде.
type StatisticsService struct {
	sessions StatisticsSessionRepo
}

func NewStatisticsService(repo StatisticsSessionRepo) *StatisticsService {
	return &StatisticsService{sessions: repo}
}

// Возвращает количество активных (не истекших) сессий.
func (s *StatisticsService) ActiveUsers(ctx context.Context) (int64, error) {
	return s.sessions.CountActive(ctx, time.Now())
}
