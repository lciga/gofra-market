package handlers

import (
	"Gofra_Market/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createListingReq struct {
	GoferName   string `json:"gofer_name" binding:"required"`
	GoferRarity int    `json:"gofer_rarity" binding:"required"`
	Price       int64  `json:"price" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type buyReq struct {
	ListingID string `json:"listing_id" binding:"required"`
}

type listingResp struct {
	ID          string           `json:"id"`
	GoferID     string           `json:"gofer_id"`
	SellerID    string           `json:"seller_id"`
	Price       int64            `json:"price"`
	IsSold      bool             `json:"is_sold"`
	BuyerID     string           `json:"buyer_id,omitempty"`
	Description string           `json:"description,omitempty"`
	Image       listingImageResp `json:"image"`
	CreatedAt   string           `json:"created_at"`
	Gofer       listingGoferResp `json:"gofer"`
}

type listingImageResp struct {
	SourceURL          *string `json:"source_url,omitempty"`
	ContentType        *string `json:"content_type,omitempty"`
	FetchedAt          *string `json:"fetched_at,omitempty"`
	DebugBase64Snippet *string `json:"debug_snippet_b64,omitempty"`
}

type listingGoferResp struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Rarity int    `json:"rarity"`
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

	id, err := h.svc.CreateWithGofer(c.Request.Context(), sellerID, req.GoferName, req.GoferRarity, req.Price, req.Description)
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

	l, g, err := h.svc.GetWithGofer(c.Request.Context(), id, requester)
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

func (h *ListingHandler) GetMyListings(c *gin.Context) {
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

	listings, gofers, err := h.svc.GetUserListingsWithGofers(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build response with full listing data including gofer info
	var result []map[string]interface{}
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

		// Show description only to seller OR buyer (if listing is sold)
		description := listing.Description
		isSeller := listing.SellerID == userID
		isBuyer := listing.IsSold && listing.BuyerID != primitive.NilObjectID && listing.BuyerID == userID
		if !isSeller && !isBuyer {
			description = ""
		}

		item := map[string]interface{}{
			"id":          listing.ID.Hex(),
			"gofer_id":    listing.GoferID.Hex(),
			"seller_id":   listing.SellerID.Hex(),
			"buyer_id":    buyerID,
			"price":       listing.Price,
			"is_sold":     listing.IsSold,
			"description": description,
			"created_at":  listing.CreatedAt.Format(time.RFC3339),
			"gofer": map[string]interface{}{
				"id":     gofers[i].ID.Hex(),
				"name":   gofers[i].Name,
				"rarity": gofers[i].Rarity,
			},
			"image": map[string]interface{}{
				"source_url":        listing.Image.SourceURL,
				"content_type":      listing.Image.ContentType,
				"fetched_at":        fetchedAt,
				"debug_snippet_b64": listing.Image.DebugSnippet,
			},
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{"listings": result})
}

func (h *ListingHandler) GetMyGofers(c *gin.Context) {
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

	gofers, err := h.svc.GetUserGofers(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]map[string]interface{}, 0, len(gofers))
	for _, gofer := range gofers {
		result = append(result, map[string]interface{}{
			"id":         gofer.ID.Hex(),
			"name":       gofer.Name,
			"rarity":     gofer.Rarity,
			"created_at": gofer.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{"gofers": result})
}
