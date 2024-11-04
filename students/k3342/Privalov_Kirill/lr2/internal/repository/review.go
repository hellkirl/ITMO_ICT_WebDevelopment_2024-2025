package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"lr2/internal/models"
	"lr2/pkg/logger"
)

type ReviewRepositoryInterface interface {
	CreateReview(ctx context.Context, review *models.Review) error
	GetReviewByID(ctx context.Context, reviewID int) (*models.Review, error)
	UpdateReview(ctx context.Context, review *models.Review) error
	DeleteReview(ctx context.Context, reviewID int) error
}

type ReviewRepository struct {
	db     *pgxpool.Pool
	logger *logger.ZapLogger
}

func NewReviewRepository(db *pgxpool.Pool, logger *logger.ZapLogger) *ReviewRepository {
	return &ReviewRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ReviewRepository) CreateReview(ctx context.Context, review *models.Review) error {
	var status string
	checkQuery := `
        SELECT status FROM reservations
        WHERE id = $1
    `
	err := r.db.QueryRow(ctx, checkQuery, review.ReservationID).Scan(&status)
	if err != nil {
		r.logger.Error("Failed to verify reservation status for review", zap.Error(err))
		return err
	}
	if status != "checked_out" {
		return errors.New("cannot write review for an active reservation")
	}

	insertQuery := `
        INSERT INTO reviews (reservation_id, rating, comment, created_at)
        VALUES ($1, $2, $3, NOW())
        RETURNING id, created_at
    `
	err = r.db.QueryRow(ctx, insertQuery, review.ReservationID, review.Rating, review.Comment).Scan(&review.ID, &review.CreatedAt)
	if err != nil {
		r.logger.Error("Failed to create review", zap.Error(err))
		return err
	}

	r.logger.Info("Review created successfully", zap.Int("review_id", review.ID))
	return nil
}

func (r *ReviewRepository) GetReviewByID(ctx context.Context, reviewID int) (*models.Review, error) {
	query := `
        SELECT id, reservation_id, rating, comment, created_at
        FROM reviews
        WHERE id = $1
    `
	review := &models.Review{}
	err := r.db.QueryRow(ctx, query, reviewID).Scan(
		&review.ID,
		&review.ReservationID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
	)
	if err != nil {
		r.logger.Error("Failed to retrieve review by ID", zap.Error(err))
		return nil, err
	}
	return review, nil
}

func (r *ReviewRepository) UpdateReview(ctx context.Context, review *models.Review) error {
	query := `
        UPDATE reviews
        SET rating = $1, comment = $2
        WHERE id = $3
    `
	cmd, err := r.db.Exec(ctx, query, review.Rating, review.Comment, review.ID)
	if err != nil {
		r.logger.Error("Failed to update review", zap.Error(err))
		return err
	}
	if cmd.RowsAffected() == 0 {
		errMsg := "review not found"
		r.logger.Warn(errMsg, zap.Int("review_id", review.ID))
		return errors.New(errMsg)
	}
	r.logger.Info("Review updated successfully", zap.Int("review_id", review.ID))
	return nil
}

func (r *ReviewRepository) DeleteReview(ctx context.Context, reviewID int) error {
	query := `
        DELETE FROM reviews
        WHERE id = $1
    `
	cmd, err := r.db.Exec(ctx, query, reviewID)
	if err != nil {
		r.logger.Error("Failed to delete review", zap.Error(err))
		return err
	}
	if cmd.RowsAffected() == 0 {
		errMsg := "review not found"
		r.logger.Warn(errMsg, zap.Int("review_id", reviewID))
		return errors.New(errMsg)
	}
	r.logger.Info("Review deleted successfully", zap.Int("review_id", reviewID))
	return nil
}
