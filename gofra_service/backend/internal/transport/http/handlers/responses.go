package handlers

import "Gofra_Market/internal/service"

// Описывает стандартную структуру ошибки API.
type errorResponse struct {
	Error string `json:"error"`
}

// Используется в swagger для описания ответа поиска маркета.
type marketSearchResponse struct {
	Items []serviceCard `json:"items"`
	Total int64         `json:"total"`
}

// Дублирует сервисный ответ, чтобы swagger не подтягивал лишние поля.
type serviceCard struct {
	ID          string          `json:"id"`
	GoferID     string          `json:"gofer_id"`
	SellerID    string          `json:"seller_id"`
	BuyerID     string          `json:"buyer_id,omitempty"`
	Price       int64           `json:"price"`
	IsSold      bool            `json:"is_sold"`
	Description string          `json:"description,omitempty"`
	CreatedAt   string          `json:"created_at"`
	Gofer       marketGoferResp `json:"gofer"`
	Image       marketImageResp `json:"image"`
}

// Повторяет структуру gofera для swagger.
type marketGoferResp struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Rarity int    `json:"rarity"`
}

// Описывает изображение в карточке.
type marketImageResp struct {
	SourceURL          *string `json:"source_url,omitempty"`
	ContentType        *string `json:"content_type,omitempty"`
	FetchedAt          *string `json:"fetched_at,omitempty"`
	DebugBase64Snippet *string `json:"debug_snippet_b64,omitempty"`
}

// Описывает ответ со списком листингов текущего пользователя.
type listingListResponse struct {
	Listings []listingResp `json:"listings"`
}

// Описывает ответ со списком гоферов пользователя.
type goferListResponse struct {
	Gofers []goferItem `json:"gofers"`
}

// Оисывает карточку гофера.
type goferItem struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Rarity    int    `json:"rarity"`
	CreatedAt string `json:"created_at"`
}

// Преобразует сервисную карточку для swagger-ответа.
func toServiceCard(c service.Card) serviceCard {
	return serviceCard{
		ID:          c.ID,
		GoferID:     c.GoferID,
		SellerID:    c.SellerID,
		BuyerID:     c.BuyerID,
		Price:       c.Price,
		IsSold:      c.IsSold,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		Gofer: marketGoferResp{
			ID:     c.Gofer.ID,
			Name:   c.Gofer.Name,
			Rarity: c.Gofer.Rarity,
		},
		Image: marketImageResp{
			SourceURL:          c.Image.SourceURL,
			ContentType:        c.Image.ContentType,
			FetchedAt:          c.Image.FetchedAt,
			DebugBase64Snippet: c.Image.DebugBase64Snippet,
		},
	}
}
