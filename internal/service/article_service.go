package service

import (
	"context"

	request_models "github.com/NarzhanProduction/Geography/internal/api/models"
	"github.com/NarzhanProduction/Geography/internal/database/repository"
	"github.com/NarzhanProduction/Geography/internal/pkg/logger"
)

type ArticleService struct {
	repo   repository.ArticleRepository
	logger logger.Logger
}

func NewArticleService(artcRepo repository.ArticleRepository, lgr logger.Logger) *ArticleService {
	return &ArticleService{
		repo:   artcRepo,
		logger: lgr,
	}
}

func (s *ArticleService) GetArticle(ctx context.Context, id int) (*request_models.GetArticleResponse, error) {
	article, err := s.repo.GetArticle(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := &request_models.GetArticleResponse{
		Title:     article.Title,
		HTML:      article.HTML,
		UpdatedAt: article.UpdatedAt,
	}
	return resp, nil
}

func (s *ArticleService) UpdateArticle(ctx context.Context, id int, html string) error {
	return s.repo.UpdateArticle(ctx, id, html)
}
