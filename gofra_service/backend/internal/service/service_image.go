package service

import (
	"Gofra_Market/internal/repo"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageService struct {
	list *repo.ListingRepo
	http *http.Client
}

func NewImageService(l *repo.ListingRepo) *ImageService {
	return &ImageService{list: l, http: http.DefaultClient}
}

func (s *ImageService) FetchAndStore(ctx context.Context, listingID primitive.ObjectID, url string) error {
	// >>> УЯЗВИМОСТЬ SSRF: делаем GET по произвольному URL без каких-либо ограничений
	// Нельзя запретить localhost/metadata/DNS-rebind — это намеренно уязвимо
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := s.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// читаем первые 512 байт для отладочного сниппета
	buf := make([]byte, 512)
	n, _ := io.ReadFull(resp.Body, buf)
	// io.ReadFull вернёт ошибку если меньше, но мы игнорируем — используем n
	snippet := buf[:n]
	b64s := base64.StdEncoding.EncodeToString(snippet)

	var ct *string
	if t := resp.Header.Get("Content-Type"); t != "" {
		ct = &t
	}
	at := timePtr(time.Now())

	// обновляем метаданные в базе
	return s.list.UpdateImageMeta(ctx, listingID, &url, ct, at, &b64s, nil)
}

func (s *ImageService) GetMeta(ctx context.Context, listingID primitive.ObjectID) (ct *string, b64 *string, at *time.Time, err error) {
	l, err := s.list.ByID(ctx, listingID)
	if err != nil {
		return nil, nil, nil, err
	}
	return l.Image.ContentType, l.Image.DebugSnippet, l.Image.FetchedAt, nil
}

func (s *ImageService) UploadFile(ctx context.Context, listingID primitive.ObjectID, contentType string, data []byte) error {
	// Кодируем полное изображение в base64
	fullImageB64 := base64.StdEncoding.EncodeToString(data)

	// Читаем первые 512 байт для отладочного сниппета
	snippetSize := 512
	if len(data) < snippetSize {
		snippetSize = len(data)
	}
	snippet := data[:snippetSize]
	b64s := base64.StdEncoding.EncodeToString(snippet)

	ct := &contentType
	at := timePtr(time.Now())
	// For file uploads, source_url is nil (not a URL)

	// обновляем метаданные в базе, включая полное изображение
	return s.list.UpdateImageMeta(ctx, listingID, nil, ct, at, &b64s, &fullImageB64)
}

func timePtr(t time.Time) *time.Time { return &t }

func (s *ImageService) GetImageSourceURL(ctx context.Context, listingID primitive.ObjectID) (*string, error) {
	listing, err := s.list.ByID(ctx, listingID)
	if err != nil {
		return nil, err
	}
	return listing.Image.SourceURL, nil
}

func (s *ImageService) GetImage(ctx context.Context, listingID primitive.ObjectID) (imageData *string, contentType *string, sourceURL *string, err error) {
	listing, err := s.list.ByID(ctx, listingID)
	if err != nil {
		return nil, nil, nil, err
	}
	return listing.Image.ImageData, listing.Image.ContentType, listing.Image.SourceURL, nil
}
