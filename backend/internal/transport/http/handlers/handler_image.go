package handlers

import (
	"Gofra_Market/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type fetchImageReq struct {
	URL string `json:"url" binding:"required"`
}

type imageMetaResp struct {
	ContentType        *string `json:"content_type"`
	FetchedAt          *string `json:"fetched_at"`
	DebugBase64Snippet *string `json:"debug_snippet_b64,omitempty"`
}

type ImageHandler struct {
	svc *service.ImageService
}

func NewImageHandler(s *service.ImageService) *ImageHandler {
	return &ImageHandler{svc: s}
}

func (h *ImageHandler) FetchFromUrl(c *gin.Context) {
	// >>> УЯЗВИМОСТЬ SSRF: auth required but no URL checks; id := Param("id"); JSON fetchImageReq;
	var req fetchImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// require authentication
	if _, ok := c.Get("userID"); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.FetchAndStore(c.Request.Context(), id, req.URL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ImageHandler) GetMeta(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ct, b64, at, err := h.svc.GetMeta(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var atStr *string
	if at != nil {
		s := at.Format(time.RFC3339)
		atStr = &s
	}

	resp := imageMetaResp{ContentType: ct, FetchedAt: atStr, DebugBase64Snippet: b64}
	c.JSON(http.StatusOK, resp)
}

func (h *ImageHandler) GetImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Try to get image metadata
	_, _, _, err = h.svc.GetMeta(c.Request.Context(), id)
	if err != nil {
		// No image found, return 404
		c.Status(http.StatusNotFound)
		return
	}

	// For now, just return 404 since we don't actually store images
	// Frontend will fallback to placeholder
	c.Status(http.StatusNotFound)
}

func (h *ImageHandler) UploadFile(c *gin.Context) {
	// require authentication
	if _, ok := c.Get("userID"); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Get uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	// Limit file size to 5MB
	const maxSize = 5 * 1024 * 1024
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 5MB)"})
		return
	}

	// Open file
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer f.Close()

	// Read file content
	data := make([]byte, file.Size)
	if _, err := f.Read(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file content"})
		return
	}

	// Store metadata
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	if err := h.svc.UploadFile(c.Request.Context(), id, contentType, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
