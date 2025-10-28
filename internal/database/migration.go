package database

import (
	"context"
	"fmt"

	"github.com/NarzhanProduction/Geography/internal/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(ctx context.Context, pool *pgxpool.Pool, logger logger.Logger) error {
	sqls := []string{
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`,

		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			is_admin BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
		);`,

		`CREATE TABLE IF NOT EXISTS articles (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			html TEXT NOT NULL,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
		);`,
	}

	for _, q := range sqls {
		if _, err := pool.Exec(ctx, q); err != nil {
			return fmt.Errorf("failed migration step: %w -- sql: %s", err, q)
		}
	}

	var count int
	err := pool.QueryRow(ctx, `SELECT COUNT(1) FROM articles WHERE id = 1`).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check default article: %w", err)
	}

	if count == 0 {
		logger.Info(ctx, "no default article founded, creating")
		_, err = pool.Exec(ctx, `
			INSERT INTO articles (id, title, html, updated_at)
			VALUES (1, $1, $2, now())
		`, "Басты мақала", "<h2>Бұл сайттағы алғашқы мақала</h2><p>Мұнда бастапқы мәтін сақталған.</p>")
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("couldn't insert default article: %v", err))
			return fmt.Errorf("failed to insert default article: %w", err)
		}
	}

	return nil
}
