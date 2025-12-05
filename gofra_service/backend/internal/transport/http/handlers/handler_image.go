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

// Структура для загрузки изображений
type fetchImageReq struct {
	URL string `json:"url" binding:"required"` // URL изображения
}

// Запрос метаданных изображения
type imageMetaResp struct {
	ContentType        *string `json:"content_type"`                // Тип содержимого
	FetchedAt          *string `json:"fetched_at"`                  // Время загрузки
	DebugBase64Snippet *string `json:"debug_snippet_b64,omitempty"` // Сниппет для отладки
}

// Структура хэндлера для изображений
type ImageHandler struct {
	svc *service.ImageService // Сервис для работы с изображениями
}

// Создание нового хэндлера
func NewImageHandler(s *service.ImageService) *ImageHandler {
	return &ImageHandler{svc: s}
}

// Метод для получения изображения
func (h *ImageHandler) FetchFromUrl(c *gin.Context) {
	var req fetchImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

// Метод для получения метаданных
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

// Метод для получения изображения
func (h *ImageHandler) GetImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	imageData, contentType, sourceURL, err := h.svc.GetImage(c.Request.Context(), id)
	if err != nil {
		logger.Errorf("GetImage: listing not found: %s, error: %v", id.Hex(), err)
		c.Status(http.StatusNotFound)
		return
	}

	logger.Infof("GetImage: listing=%s hasSourceURL=%v hasImageData=%v", id.Hex(), sourceURL != nil, imageData != nil)

	if sourceURL != nil && *sourceURL != "" {
		logger.Infof("GetImage: redirecting to source URL: %s", *sourceURL)
		c.Redirect(http.StatusFound, *sourceURL)
		return
	}

	if imageData != nil && *imageData != "" {
		logger.Infof("GetImage: returning image data, size: %d", len(*imageData))
		decoded, err := base64.StdEncoding.DecodeString(*imageData)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		if contentType != nil && *contentType != "" {
			c.Header("Content-Type", *contentType)
		} else {
			c.Header("Content-Type", "application/octet-stream")
		}

		c.Data(http.StatusOK, c.GetHeader("Content-Type"), decoded)
		return
	}

	c.Status(http.StatusNotFound)
}

// Метод для загрузки файлов
func (h *ImageHandler) UploadFile(c *gin.Context) {
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

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	const maxSize = 5 * 1024 * 1024
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 5MB)"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer f.Close()

	data := make([]byte, file.Size)
	if _, err := f.Read(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file content"})
		return
	}

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	if err := h.svc.UploadFile(c.Request.Context(), id, contentType, data); err != nil {
		logger.Errorf("UploadFile: failed to upload image for listing %s: %v", id.Hex(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Infof("Successfully uploaded image for listing: %s, content-type: %s, size: %d", id.Hex(), contentType, len(data))

	c.Status(http.StatusNoContent)
}
