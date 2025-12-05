// Пакет для работы с хэндлерами
package handlers

import (
	"Gofra_Market/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Структура запроса на регистрацию
type registerReq struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Структура запроса на вход
type loginReq struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Структура ответа на запрос текущего пользователя
type meResp struct {
	UserID  string `json:"user_id"`
	Login   string `json:"login"`
	Balance int64  `json:"balance"`
}

// Структура хэндлера аутентификации
type AuthHandler struct {
	svc    *service.AuthService // Сервис аутентификации
	cookie string               // Куки
}

// Создание нового хэндлера
func NewAuthHandler(s *service.AuthService, c string) *AuthHandler {
	return &AuthHandler{svc: s, cookie: c}
}

// Метод для регистрации
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, sid, err := h.svc.Register(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(h.cookie, sid, 60*60*24*30, "/", "", false, true)
	c.JSON(http.StatusCreated, meResp{UserID: user.ID.Hex(), Login: user.Login, Balance: user.Balance})
}

// Метод для входа
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, sid, err := h.svc.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie(h.cookie, sid, 60*60*24*30, "/", "", false, true)
	c.JSON(http.StatusOK, meResp{UserID: user.ID.Hex(), Login: user.Login, Balance: user.Balance})
}

// Метод получения текущего пользователя
func (h *AuthHandler) Me(c *gin.Context) {
	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	uid, ok := v.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id in context"})
		return
	}

	user, err := h.svc.Me(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, meResp{UserID: user.ID.Hex(), Login: user.Login, Balance: user.Balance})
}
