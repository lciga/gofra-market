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

// Search godoc
// @Summary Поиск карточек в маркетплейсе
// @Description Возвращает список карточек гоферов по фильтру. Уязвим для NoSQL-инъекций.
// @Tags market
// @Accept json
// @Produce json
// @Param filter query string false "MongoDB JSON фильтр"
// @Param limit query int false "Количество элементов на страницу"
// @Param page query int false "Номер страницы"
// @Param sort query string false "Поле сортировки (price_asc, price_desc)"
// @Success 200 {object} marketSearchResponse
// @Failure 500 {object} errorResponse
// @Router /api/market [get]
func (h *MarketHandler) Search(c *gin.Context) {
	raw := c.Query("filter")
	limitStr := c.DefaultQuery("limit", "20")
	pageStr := c.DefaultQuery("page", "1")
	sort := c.DefaultQuery("sort", "")

	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)

	items, total, err := h.svc.SearchRaw(c.Request.Context(), raw, limit, page, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	converted := make([]serviceCard, 0, len(items))
	for _, item := range items {
		converted = append(converted, toServiceCard(item))
	}

	c.JSON(http.StatusOK, marketSearchResponse{Items: converted, Total: total})
}
