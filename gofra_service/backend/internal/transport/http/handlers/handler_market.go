package handlers

import (
	"Gofra_Market/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Структура хэндлера для маркета
type MarketHandler struct {
	svc *service.MarketService // Сервис для работы с маркетом
}

// Создание нового хэндлера
func NewMarketHandler(s *service.MarketService) *MarketHandler {
	return &MarketHandler{svc: s}
}

// Поиск карточек (уязвим к NoSQL-инъекциям)
func (h *MarketHandler) Search(c *gin.Context) {
	raw := c.Query("filter")
	limitStr := c.DefaultQuery("limit", "20")
	pageStr := c.DefaultQuery("page", "1")
	sort := c.DefaultQuery("sort", "")

	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)

	items, total, err := h.svc.SearchRaw(c.Request.Context(), raw, limit, page, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": total})
}
