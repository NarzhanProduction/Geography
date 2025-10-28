package repository

import (
	"context"
	"fmt"

	"github.com/NarzhanProduction/Geography/internal/database/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArticleRepository interface {
	GetArticle(ctx context.Context, id int) (*models.Article, error)
	UpdateArticle(ctx context.Context, id int, html string) error
}

type articleRepository struct {
	db *pgxpool.Pool
}

func NewArticleRepository(db *pgxpool.Pool) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) GetArticle(ctx context.Context, id int) (*models.Article, error) {
	article := &models.Article{}
	query := `SELECT id, title, html, updated_at FROM articles WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&article.ID, &article.Title, &article.HTML, &article.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("article not found: %w", err)
	}
	return article, nil
}

func (r *articleRepository) UpdateArticle(ctx context.Context, id int, html string) error {
	query := `UPDATE articles SET html = $1, updated_at = NOW()
		WHERE id = $2`
	_, err := r.db.Exec(ctx, query, html, id)
	if err != nil {
		return fmt.Errorf("error updating article: %w", err)
	}
	return nil
}
