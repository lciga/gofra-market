package handlers

import (
	"Gofra_Market/internal/logger"
	"Gofra_Market/internal/service"
	"encoding/base64"
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

	// Get image data from service
	imageData, contentType, sourceURL, err := h.svc.GetImage(c.Request.Context(), id)
	if err != nil {
		// No image found, return 404
		logger.Errorf("GetImage: listing not found: %s, error: %v", id.Hex(), err)
		c.Status(http.StatusNotFound)
		return
	}

	logger.Infof("GetImage: listing=%s hasSourceURL=%v hasImageData=%v", id.Hex(), sourceURL != nil, imageData != nil)

	// If we have a source URL (from URL upload), redirect to it
	if sourceURL != nil && *sourceURL != "" {
		logger.Infof("GetImage: redirecting to source URL: %s", *sourceURL)
		c.Redirect(http.StatusFound, *sourceURL)
		return
	}

	// For file uploads, return the base64 decoded image data
	if imageData != nil && *imageData != "" {
		logger.Infof("GetImage: returning image data, size: %d", len(*imageData))
		// Decode base64
		decoded, err := base64.StdEncoding.DecodeString(*imageData)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		// Set content type
		if contentType != nil && *contentType != "" {
			c.Header("Content-Type", *contentType)
		} else {
			c.Header("Content-Type", "application/octet-stream")
		}

		// Return image data
		c.Data(http.StatusOK, c.GetHeader("Content-Type"), decoded)
		return
	}

	// No image data available
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
		logger.Errorf("UploadFile: failed to upload image for listing %s: %v", id.Hex(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log successful upload
	logger.Infof("Successfully uploaded image for listing: %s, content-type: %s, size: %d", id.Hex(), contentType, len(data))

	c.Status(http.StatusNoContent)
}
