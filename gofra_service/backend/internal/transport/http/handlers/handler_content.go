package handlers

import (
	"net/http"

	"Gofra_Market/internal/service"

	"github.com/gin-gonic/gin"
)

// Структура запроса материала на проверку
type contentReviewReq struct {
	Title string `json:"title" binding:"required"` // Заголовок материала
	Text  string `json:"text" binding:"required"`  // Текст материала
}

// Структура хэндлера материалов
type ContentHandler struct {
	svc *service.ContentService // Сервис материалов
}

// Создание нового хэндлера материалов
func NewContentHandler(s *service.ContentService) *ContentHandler {
	return &ContentHandler{svc: s}
}

// Метод подготовки письма редактору
func (h *ContentHandler) Submit(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}

	var req contentReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	review, err := h.svc.PrepareReview(c.Request.Context(), uid, req.Title, req.Text)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, errorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}
