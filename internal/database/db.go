package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

func New(ctx context.Context, dsn string) (*Postgres, error) {
	pgOnce.Do(func() {
		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			panic(fmt.Errorf("unable to parse DB config: %w", err))
		}
		// Настройка пула (опционально)
		config.MaxConns = 10 // Из конфига
		config.MinConns = 0
		config.MaxConnIdleTime = 30 * time.Minute // Закрывать idle-соединения

		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			panic(fmt.Errorf("unable to create connection pool: %w", err))
		}
		// Пинг для проверки
		if err := pool.Ping(ctx); err != nil {
			panic(fmt.Errorf("unable to ping DB: %w", err))
		}

		pgInstance = &Postgres{DB: pool}
	})
	return pgInstance, nil
}

func (pg *Postgres) Close() {
	if pgInstance != nil {
		pgInstance.DB.Close()
	}
}
