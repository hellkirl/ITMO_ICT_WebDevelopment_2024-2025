package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"lr2/internal/models"
	"lr2/pkg/logger"
)

type UserRepositoryInterface interface {
	RegisterUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type UserRepository struct {
	db     *pgxpool.Pool
	logger *logger.ZapLogger
}

func NewUserRepository(db *pgxpool.Pool, logger *logger.ZapLogger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) RegisterUser(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (username, email, password, created_at)
        VALUES ($1, $2, $3, NOW())
        RETURNING id, created_at
    `
	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		r.logger.Error("Couldn't register new user", zap.Error(err))
		return err
	}
	r.logger.Info("User is registered", zap.Int("user_id", user.ID))
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	query := `
        SELECT id, username, email, password, created_at
        FROM users
        WHERE id = $1
    `
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		r.logger.Error("Couldn't find user by id", zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
        SELECT id, username, email, password, created_at
        FROM users
        WHERE email = $1
    `
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		r.logger.Error("Couldn't find user by email", zap.Error(err))
		return nil, err
	}
	return user, nil
}
