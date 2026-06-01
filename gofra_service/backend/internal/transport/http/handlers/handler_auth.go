package handlers

import (
	"Gofra_Market/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Структура запроса на регистрацию
type registerReq struct {
	Login    string `json:"login" binding:"required"`    // Логин
	Password string `json:"password" binding:"required"` // Пароль
}

// Структура запроса на вход
type loginReq struct {
	Login    string `json:"login" binding:"required"`    // Логин
	Password string `json:"password" binding:"required"` // Пароль
}

// Структура ответа текущего пользователя
type meResp struct {
	UserID  string `json:"user_id"`         // Идентификатор пользователя
	Login   string `json:"login"`           // Логин
	Role    string `json:"role"`            // Роль
	Balance int64  `json:"balance"`         // Баланс
	Token   string `json:"token,omitempty"` // SID сессии
}

// Структура хэндлера аутентификации
type AuthHandler struct {
	svc    *service.AuthService // Сервис аутентификации
	cookie string               // Имя cookie сессии
}

// Создание нового хэндлера аутентификации
func NewAuthHandler(s *service.AuthService, c string) *AuthHandler {
	return &AuthHandler{svc: s, cookie: c}
}

// Метод регистрации пользователя
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
	c.JSON(http.StatusCreated, mapMeResp(user.ID.Hex(), user.Login, user.Role, user.Balance, sid))
}

// Метод входа пользователя
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
	c.JSON(http.StatusOK, mapMeResp(user.ID.Hex(), user.Login, user.Role, user.Balance, sid))
}

// Метод получения текущего пользователя
func (h *AuthHandler) Me(c *gin.Context) {
	uid, err := currentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: err.Error()})
		return
	}

	user, err := h.svc.Me(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, mapMeResp(user.ID.Hex(), user.Login, user.Role, user.Balance, ""))
}

// Преобразование пользователя в ответ API
func mapMeResp(id, login, role string, balance int64, token string) meResp {
	return meResp{
		UserID:  id,
		Login:   login,
		Role:    role,
		Balance: balance,
		Token:   token,
	}
}
