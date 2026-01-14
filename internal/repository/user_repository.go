package repository

import (
	"base-service/internal/model"
	"context"
	"database/sql"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (username, email, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id
	`, user.Username, user.Email)
	return err
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	return user, err
}
