package handlers

import (
	"Gofra_Market/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createListingReq struct {
	GoferID     string `json:"gofer_id" binding:"required"`
	Price       int64  `json:"price" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type buyReq struct {
	ListingID string `json:"listing_id" binding:"required"`
}

type listingResp struct {
	ID          string `json:"id"`
	GoferID     string `json:"gofer_id"`
	SellerID    string `json:"seller_id"`
	Price       int64  `json:"price"`
	IsSold      bool   `json:"is_sold"`
	BuyerID     string `json:"buyer_id,omitempty"`
	Description string `json:"description,omitempty"`
	Image       struct {
		ContentType        *string `json:"content_type"`
		FetchedAt          *string `json:"fetched_at"`
		DebugBase64Snippet *string `json:"debug_snippet_b64,omitempty"`
	} `json:"image"`
	CreatedAt string `json:"created_at"`
}

type ListingHandler struct {
	svc *service.ListingService
}

func NewListingHandler(s *service.ListingService) *ListingHandler {
	return &ListingHandler{svc: s}
}

func (h *ListingHandler) Create(c *gin.Context) {
	var req createListingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	sellerID, ok := v.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	goferID, err := primitive.ObjectIDFromHex(req.GoferID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid gofer id"})
		return
	}

	id, err := h.svc.Create(c.Request.Context(), sellerID, goferID, req.Price, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id.Hex()})
}

func (h *ListingHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var requester *primitive.ObjectID
	if v, ok := c.Get("userID"); ok {
		if uid, ok := v.(primitive.ObjectID); ok {
			requester = &uid
		}
	}

	l, err := h.svc.Get(c.Request.Context(), id, requester)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
	}
	if l.Image.ContentType != nil {
		resp.Image.ContentType = l.Image.ContentType
	}
	resp.Image.FetchedAt = fetchedAt
	if l.Image.DebugSnippet != nil {
		resp.Image.DebugBase64Snippet = l.Image.DebugSnippet
	}

	// If requester is not seller or buyer, service.Get already zeroed Description
	c.JSON(http.StatusOK, resp)
}

func (h *ListingHandler) Buy(c *gin.Context) {
	var req buyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	buyerID, ok := v.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	listingID, err := primitive.ObjectIDFromHex(req.ListingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid listing id"})
		return
	}

	if err := h.svc.Buy(c.Request.Context(), buyerID, listingID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ListingHandler) Bump(c *gin.Context) {
	v, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	userID, ok := v.(primitive.ObjectID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	idStr := c.Param("id")
	listingID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.Bump(c.Request.Context(), userID, listingID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
