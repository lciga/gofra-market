package handlers

import (
	"net/http"

	"Gofra_Market/internal/service"

	"github.com/gin-gonic/gin"
)

// Структура хэндлера статистики
type StatisticsHandler struct {
	statistics *service.StatisticsService // Сервис статистики
}

// Создание нового хэндлера статистики
func NewStatisticsHandler(statSvc *service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{statistics: statSvc}
}

// Метод получения количества активных пользователей
func (h *StatisticsHandler) ActiveUsers(c *gin.Context) {
	count, err := h.statistics.ActiveUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"active_users": count})
}

// Метод получения статистики посещений
func (h *StatisticsHandler) Visits(c *gin.Context) {
	stats, err := h.statistics.Visits(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
