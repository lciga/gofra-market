package app

import (
	"Gofra_Market/internal/transport/http/handlers"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth    *handlers.AuthHandler
	Market  *handlers.MarketHandler
	Listing *handlers.ListingHandler
	Image   *handlers.ImageHandler
}

func RegisterRoutes(e *gin.Engine, h Handlers) {
	// public
	e.POST("/api/register", h.Auth.Register)
	e.POST("/api/login", h.Auth.Login)
	e.GET("/api/market", h.Market.Search)
	e.GET("/api/listings/:id", h.Listing.Get)

	// protected (middleware should set userID in context)
	api := e.Group("/api")
	api.GET("/me", h.Auth.Me)
	api.POST("/listings", h.Listing.Create)
	api.POST("/buy", h.Listing.Buy)
	api.POST("/listings/:id/bump", h.Listing.Bump)
	api.POST("/listings/:id/image_from_url", h.Image.FetchFromUrl)
	api.GET("/listings/:id/image/meta", h.Image.GetMeta)
}
