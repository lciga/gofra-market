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

// Register godoc
// @Summary Регистрация пользователя
// @Description Создаёт пользователя и авторизует его через cookie.
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body registerReq true "Данные пользователя"
// @Success 201 {object} meResp
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	user, sid, err := h.svc.Register(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.SetCookie(h.cookie, sid, 60*60*24*30, "/", "", false, true)
	c.JSON(http.StatusCreated, meResp{UserID: user.ID.Hex(), Login: user.Login, Balance: user.Balance})
}

// Login godoc
// @Summary Авторизация пользователя
// @Description Аутентифицирует пользователя и устанавливает cookie сессии.
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body loginReq true "Учётные данные"
// @Success 200 {object} meResp
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	user, sid, err := h.svc.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}

	c.SetCookie(h.cookie, sid, 60*60*24*30, "/", "", false, true)
	c.JSON(http.StatusOK, meResp{UserID: user.ID.Hex(), Login: user.Login, Balance: user.Balance})
}

// Me godoc
// @Summary Текущий пользователь
// @Description Возвращает текущего авторизованного пользователя.
// @Tags auth
// @Produce json
// @Success 200 {object} meResp
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /api/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthenticated"})
		return
	}
	uid, ok := v.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "invalid user id in context"})
		return
	}

	user, err := h.svc.Me(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, meResp{UserID: user.ID.Hex(), Login: user.Login, Balance: user.Balance})
}
