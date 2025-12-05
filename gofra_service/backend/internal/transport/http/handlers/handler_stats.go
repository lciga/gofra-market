package handlers

import (
	"net/http"

	"Gofra_Market/internal/service"

	"github.com/gin-gonic/gin"
)

// Обслуживает эндпоинты статистики по пользователям.
type StatisticsHandler struct {
	statistics *service.StatisticsService
}

func NewStatisticsHandler(statSvc *service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{statistics: statSvc}
}

// Возвращает количество активных пользователей на сервере.
func (h *StatisticsHandler) ActiveUsers(c *gin.Context) {
	count, err := h.statistics.ActiveUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"active_users": count})
}
