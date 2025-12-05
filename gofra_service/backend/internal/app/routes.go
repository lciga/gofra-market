// Пакет для работы сервера и роутера
package app

import (
	"Gofra_Market/internal/transport/http/handlers"

	"github.com/gin-gonic/gin"
)

// Структура для хранения хэндлеров
type Handlers struct {
	Auth    *handlers.AuthHandler       // Хэндлер аутентификации
	Market  *handlers.MarketHandler     // Хэндлер маркета
	Listing *handlers.ListingHandler    // Хэндлер листинга
	Image   *handlers.ImageHandler      // Хэндлер для работы с изображениями
	Stats   *handlers.StatisticsHandler // Хэндлер статистики
}

// Регистрация роутов
func RegisterRoutes(e *gin.Engine, h Handlers) {
	e.POST("/api/register", h.Auth.Register)
	e.POST("/api/login", h.Auth.Login)
	e.GET("/api/market", h.Market.Search)
	e.GET("/api/listings/:id", h.Listing.Get)
	e.GET("/api/listings/:id/image", h.Image.GetImage)

	api := e.Group("/api")
	api.GET("/me", h.Auth.Me)
	api.GET("/my-listings", h.Listing.GetMyListings)
	api.GET("/my-gofers", h.Listing.GetMyGofers)
	api.POST("/listings", h.Listing.Create)
	api.POST("/buy", h.Listing.Buy)
	api.POST("/listings/:id/bump", h.Listing.Bump)
	api.POST("/listings/:id/image_from_url", h.Image.FetchFromUrl)
	api.POST("/listings/:id/image_upload", h.Image.UploadFile)
	api.GET("/listings/:id/image/meta", h.Image.GetMeta)
	api.GET("/stats/active-users", h.Stats.ActiveUsers)
}
