package handlers

import (
	"Gofra_Market/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Запрос на создание листинга
type createListingReq struct {
	GoferName   string `json:"gofer_name" binding:"required"`   // Имя гофера
	GoferRarity int    `json:"gofer_rarity" binding:"required"` // Редкость
	Price       int64  `json:"price" binding:"required"`        // Цена
	Description string `json:"description" binding:"required"`  // Описание
}

// Ответ на создание листинга
type createListingResp struct {
	ID string `json:"id"`
}

// Запрос на покупку
type buyReq struct {
	ListingID string `json:"listing_id" binding:"required"` // Идентификатор листинга
}

// Ответ на запрос листинга
type listingResp struct {
	ID          string           `json:"id"`                    // Идентификатор листинга
	GoferID     string           `json:"gofer_id"`              // Идентификатор гофера
	SellerID    string           `json:"seller_id"`             // Идентификатор продавца
	Price       int64            `json:"price"`                 // Цена
	IsSold      bool             `json:"is_sold"`               // Флаг продажи
	BuyerID     string           `json:"buyer_id,omitempty"`    // Идентификатор покупателя
	Description string           `json:"description,omitempty"` // Описание
	Image       listingImageResp `json:"image"`                 // Изображение
	CreatedAt   string           `json:"created_at"`            // Время создания
	Gofer       listingGoferResp `json:"gofer"`                 // Гофер
}

// Ответ на запрос изображения листинга
type listingImageResp struct {
	SourceURL          *string `json:"source_url,omitempty"`        // URL источника
	ContentType        *string `json:"content_type,omitempty"`      // Тип содержимого
	FetchedAt          *string `json:"fetched_at,omitempty"`        // Время загрузки
	DebugBase64Snippet *string `json:"debug_snippet_b64,omitempty"` // Сниппет для отладки
}

// Запрос гофреа
type listingGoferResp struct {
	ID     string `json:"id"`     // Идентификатор
	Name   string `json:"name"`   // Имя
	Rarity int    `json:"rarity"` // Редкость
}

// Структура хэндлера листинга
type ListingHandler struct {
	svc *service.ListingService // Сервис для работы с листингом
}

// Создание нового хэндлера
func NewListingHandler(s *service.ListingService) *ListingHandler {
	return &ListingHandler{svc: s}
}

// Create godoc
// @Summary Создание листинга с гофером
// @Description Создаёт листинг и нового гофера, привязанного к продавцу.
// @Tags listings
// @Accept json
// @Produce json
// @Param payload body createListingReq true "Описание листинга"
// @Success 201 {object} createListingResp
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/listings [post]
func (h *ListingHandler) Create(c *gin.Context) {
	var req createListingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthenticated"})
		return
	}
	sellerID, ok := v.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "invalid user id"})
		return
	}

	id, err := h.svc.CreateWithGofer(c.Request.Context(), sellerID, req.GoferName, req.GoferRarity, req.Price, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createListingResp{ID: id.Hex()})
}

// Get godoc
// @Summary Получение листинга
// @Description Возвращает листинг с данными гофера. Скрывает описание для посторонних пользователей.
// @Tags listings
// @Produce json
// @Param id path string true "ID листинга"
// @Success 200 {object} listingResp
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /api/listings/{id} [get]
func (h *ListingHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}

	var requester *primitive.ObjectID
	if v, ok := c.Get("userID"); ok {
		if uid, ok := v.(primitive.ObjectID); ok {
			requester = &uid
		}
	}

	l, g, err := h.svc.GetWithGofer(c.Request.Context(), id, requester)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse{Error: err.Error()})
		return
	}

	var buyerID string
	if l.BuyerID != primitive.NilObjectID {
		buyerID = l.BuyerID.Hex()
	}

	var fetchedAt *string
	if l.Image.FetchedAt != nil {
		s := l.Image.FetchedAt.Format(time.RFC3339)
		fetchedAt = &s
	}

	resp := listingResp{
		ID:          l.ID.Hex(),
		GoferID:     l.GoferID.Hex(),
		SellerID:    l.SellerID.Hex(),
		Price:       l.Price,
		IsSold:      l.IsSold,
		BuyerID:     buyerID,
		Description: l.Description,
		CreatedAt:   l.CreatedAt.Format(time.RFC3339),
		Gofer: listingGoferResp{
			ID:     g.ID.Hex(),
			Name:   g.Name,
			Rarity: g.Rarity,
		},
	}
	resp.Image.SourceURL = l.Image.SourceURL
	resp.Image.ContentType = l.Image.ContentType
	resp.Image.FetchedAt = fetchedAt
	resp.Image.DebugBase64Snippet = l.Image.DebugSnippet

	c.JSON(http.StatusOK, resp)
}

// Buy godoc
// @Summary Покупка листинга
// @Description Переводит гофера покупателю и списывает средства.
// @Tags listings
// @Accept json
// @Param payload body buyReq true "Листинг для покупки"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/buy [post]
func (h *ListingHandler) Buy(c *gin.Context) {
	var req buyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthenticated"})
		return
	}
	buyerID, ok := v.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "invalid user id"})
		return
	}

	listingID, err := primitive.ObjectIDFromHex(req.ListingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid listing id"})
		return
	}

	if err := h.svc.Buy(c.Request.Context(), buyerID, listingID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// Bump godoc
// @Summary Поднять листинг
// @Description Списывает стоимость поднятия с баланса продавца (уязвим к uint underflow).
// @Tags listings
// @Param id path string true "ID листинга"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/listings/{id}/bump [post]
func (h *ListingHandler) Bump(c *gin.Context) {
	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthenticated"})
		return
	}
	userID, ok := v.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "invalid user id"})
		return
	}

	idStr := c.Param("id")
	listingID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return
	}

	if err := h.svc.Bump(c.Request.Context(), userID, listingID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetMyListings godoc
// @Summary Мои листинги
// @Description Возвращает листинги текущего пользователя вместе с гоферами.
// @Tags listings
// @Produce json
// @Success 200 {object} listingListResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/my-listings [get]
func (h *ListingHandler) GetMyListings(c *gin.Context) {
	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthenticated"})
		return
	}
	userID, ok := v.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "invalid user id"})
		return
	}

	listings, gofers, err := h.svc.GetUserListingsWithGofers(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	result := make([]listingResp, 0, len(listings))
	for i, listing := range listings {
		var buyerID string
		if listing.BuyerID != primitive.NilObjectID {
			buyerID = listing.BuyerID.Hex()
		}

		var fetchedAt *string
		if listing.Image.FetchedAt != nil {
			s := listing.Image.FetchedAt.Format(time.RFC3339)
			fetchedAt = &s
		}

		description := listing.Description
		isSeller := listing.SellerID == userID
		isBuyer := listing.IsSold && listing.BuyerID != primitive.NilObjectID && listing.BuyerID == userID
		if !isSeller && !isBuyer {
			description = ""
		}

		item := listingResp{
			ID:          listing.ID.Hex(),
			GoferID:     listing.GoferID.Hex(),
			SellerID:    listing.SellerID.Hex(),
			Price:       listing.Price,
			IsSold:      listing.IsSold,
			BuyerID:     buyerID,
			Description: description,
			CreatedAt:   listing.CreatedAt.Format(time.RFC3339),
			Gofer: listingGoferResp{
				ID:     gofers[i].ID.Hex(),
				Name:   gofers[i].Name,
				Rarity: gofers[i].Rarity,
			},
		}
		item.Image.SourceURL = listing.Image.SourceURL
		item.Image.ContentType = listing.Image.ContentType
		item.Image.FetchedAt = fetchedAt
		item.Image.DebugBase64Snippet = listing.Image.DebugSnippet

		result = append(result, item)
	}

	c.JSON(http.StatusOK, listingListResponse{Listings: result})
}

// GetMyGofers godoc
// @Summary Мои гоферы
// @Description Возвращает гоферов, которыми владеет текущий пользователь.
// @Tags listings
// @Produce json
// @Success 200 {object} goferListResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/my-gofers [get]
func (h *ListingHandler) GetMyGofers(c *gin.Context) {
	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, errorResponse{Error: "unauthenticated"})
		return
	}
	userID, ok := v.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "invalid user id"})
		return
	}

	gofers, err := h.svc.GetUserGofers(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	result := make([]goferItem, 0, len(gofers))
	for _, gofer := range gofers {
		result = append(result, goferItem{
			ID:        gofer.ID.Hex(),
			Name:      gofer.Name,
			Rarity:    gofer.Rarity,
			CreatedAt: gofer.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, goferListResponse{Gofers: result})
}
