package repository

import (
	"context"
	"fmt"

	"github.com/NarzhanProduction/Geography/internal/database/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByLogin(ctx context.Context, name string) (*models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`
	_, err := r.db.Exec(ctx, query, user.Name, user.Email)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, uuid uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, email, is_admin FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, uuid).Scan(&user.ID, &user.Name, &user.Email, &user.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, password_hash, is_admin FROM users WHERE name = $1`
	err := r.db.QueryRow(ctx, query, login).Scan(&user.ID, &user.Name, &user.PasswordHash, &user.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}
