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
	return s.list.UpdateImageMeta(ctx, listingID, url, ct, at, &b64s)
}

func (s *ImageService) GetMeta(ctx context.Context, listingID primitive.ObjectID) (ct *string, b64 *string, at *time.Time, err error) {
	l, err := s.list.ByID(ctx, listingID)
	if err != nil {
		return nil, nil, nil, err
	}
	return l.Image.ContentType, l.Image.DebugSnippet, l.Image.FetchedAt, nil
}

func timePtr(t time.Time) *time.Time { return &t }
