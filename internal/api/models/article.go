package request_models

import "time"

type ArticleUpdateRequest struct {
	ID   int    `json:"id"`
	HTML string `json:"html"`
}

type GetArticleRequest struct {
	ID int `json:"id"`
}

type GetArticleResponse struct {
	Title     string    `json:"title"`
	HTML      string    `json:"html"`
	UpdatedAt time.Time `json:"updated_at"`
}
