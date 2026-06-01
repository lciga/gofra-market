package handlers

import (
	"net/http"

	"Gofra_Market/internal/service"

	"github.com/gin-gonic/gin"
)

// Структура хэндлера админской панели
type AdminHandler struct {
	svc *service.AdminService // Сервис админской панели
}

// Создание нового хэндлера админской панели
func NewAdminHandler(s *service.AdminService) *AdminHandler {
	return &AdminHandler{svc: s}
}

// Метод получения данных админской панели
func (h *AdminHandler) Dashboard(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}

	dashboard, err := h.svc.Dashboard(c.Request.Context(), uid)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, errorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}
