package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"lr2/internal/models"
	"lr2/pkg/logger"
	"time"
)

type ReservationRepositoryInterface interface {
	CreateReservation(ctx context.Context, reservation *models.Reservation) error
	GetReservationByID(ctx context.Context, reservationID int) (*models.Reservation, error)
	UpdateReservation(ctx context.Context, reservation *models.Reservation) error
	DeleteReservation(ctx context.Context, reservationID int) error
	CheckIn(ctx context.Context, reservationID int) error
	CheckOut(ctx context.Context, reservationID int) error
	GetRecentGuests(ctx context.Context) ([]models.GuestInfo, error)
}

type ReservationRepository struct {
	db     *pgxpool.Pool
	logger *logger.ZapLogger
}

func NewReservationRepository(db *pgxpool.Pool, logger *logger.ZapLogger) *ReservationRepository {
	return &ReservationRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ReservationRepository) CreateReservation(ctx context.Context, reservation *models.Reservation) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.logger.Error("Failed to begin transaction", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	var exists bool
	checkQuery := `
        SELECT EXISTS (
            SELECT 1 FROM reservations
            WHERE room_id = $1
            AND status IN ('reserved', 'checked_in')
            AND ($2 < check_out AND $3 > check_in)
        )
    `
	err = tx.QueryRow(ctx, checkQuery, reservation.RoomID, reservation.CheckIn, reservation.CheckOut).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check room availability", zap.Error(err))
		return err
	}
	if exists {
		return errors.New("room is not available for the selected period")
	}

	insertQuery := `
        INSERT INTO reservations (user_id, room_id, check_in, check_out, status, created_at)
        VALUES ($1, $2, $3, $4, 'reserved', NOW())
        RETURNING id, created_at
    `
	err = tx.QueryRow(ctx, insertQuery, reservation.UserID, reservation.RoomID, reservation.CheckIn, reservation.CheckOut).Scan(&reservation.ID, &reservation.CreatedAt)
	if err != nil {
		r.logger.Error("Failed to create reservation", zap.Error(err))
		return err
	}

	updateRoomQuery := `
        UPDATE rooms
        SET is_available = FALSE
        WHERE id = $1
    `
	_, err = tx.Exec(ctx, updateRoomQuery, reservation.RoomID)
	if err != nil {
		r.logger.Error("Failed to update room availability", zap.Error(err))
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		r.logger.Error("Failed to commit transaction", zap.Error(err))
		return err
	}

	r.logger.Info("Reservation created successfully", zap.Int("reservation_id", reservation.ID))
	return nil
}

func (r *ReservationRepository) GetReservationByID(ctx context.Context, reservationID int) (*models.Reservation, error) {
	query := `
        SELECT id, user_id, room_id, check_in, check_out, status, created_at
        FROM reservations
        WHERE id = $1
    `
	reservation := &models.Reservation{}
	err := r.db.QueryRow(ctx, query, reservationID).Scan(
		&reservation.ID,
		&reservation.UserID,
		&reservation.RoomID,
		&reservation.CheckIn,
		&reservation.CheckOut,
		&reservation.Status,
		&reservation.CreatedAt,
	)
	if err != nil {
		r.logger.Error("Failed to get reservation by ID", zap.Error(err))
		return nil, err
	}
	return reservation, nil
}

func (r *ReservationRepository) UpdateReservation(ctx context.Context, reservation *models.Reservation) error {
	query := `
        UPDATE reservations
        SET check_in = $1, check_out = $2, status = $3
        WHERE id = $4
    `
	cmd, err := r.db.Exec(ctx, query, reservation.CheckIn, reservation.CheckOut, reservation.Status, reservation.ID)
	if err != nil {
		r.logger.Error("Failed to update reservation", zap.Error(err))
		return err
	}
	if cmd.RowsAffected() == 0 {
		errMsg := "reservation not found"
		r.logger.Warn(errMsg, zap.Int("reservation_id", reservation.ID))
		return errors.New(errMsg)
	}
	r.logger.Info("Reservation updated successfully", zap.Int("reservation_id", reservation.ID))
	return nil
}

func (r *ReservationRepository) DeleteReservation(ctx context.Context, reservationID int) error {
	query := `
        DELETE FROM reservations
        WHERE id = $1
    `
	cmd, err := r.db.Exec(ctx, query, reservationID)
	if err != nil {
		r.logger.Error("Failed to delete reservation", zap.Error(err))
		return err
	}
	if cmd.RowsAffected() == 0 {
		errMsg := "reservation not found"
		r.logger.Warn(errMsg, zap.Int("reservation_id", reservationID))
		return errors.New(errMsg)
	}
	r.logger.Info("Reservation deleted successfully", zap.Int("reservation_id", reservationID))
	return nil
}

func (r *ReservationRepository) CheckIn(ctx context.Context, reservationID int) error {
	query := `
        UPDATE reservations
        SET status = 'checked_in'
        WHERE id = $1 AND status = 'reserved'
    `
	cmd, err := r.db.Exec(ctx, query, reservationID)
	if err != nil {
		r.logger.Error("Failed to check in user", zap.Error(err))
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("reservation not found or already checked in")
	}
	r.logger.Info("User checked in successfully", zap.Int("reservation_id", reservationID))
	return nil
}

func (r *ReservationRepository) CheckOut(ctx context.Context, reservationID int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.logger.Error("Failed to begin transaction for check out", zap.Error(err))
		return err
	}
	defer tx.Rollback(ctx)

	updateReservation := `
        UPDATE reservations
        SET status = 'checked_out'
        WHERE id = $1 AND status = 'checked_in'
    `
	cmd, err := tx.Exec(ctx, updateReservation, reservationID)
	if err != nil {
		r.logger.Error("Failed to update reservation status to checked_out", zap.Error(err))
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("reservation not found or not checked in")
	}

	var roomID int
	getRoomQuery := `
        SELECT room_id FROM reservations
        WHERE id = $1
    `
	err = tx.QueryRow(ctx, getRoomQuery, reservationID).Scan(&roomID)
	if err != nil {
		r.logger.Error("Failed to get room_id from reservation", zap.Error(err))
		return err
	}

	updateRoom := `
        UPDATE rooms
        SET is_available = TRUE
        WHERE id = $1
    `
	_, err = tx.Exec(ctx, updateRoom, roomID)
	if err != nil {
		r.logger.Error("Failed to update room availability on check out", zap.Error(err))
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		r.logger.Error("Failed to commit transaction for check out", zap.Error(err))
		return err
	}

	r.logger.Info("User checked out successfully", zap.Int("reservation_id", reservationID))
	return nil
}

func (r *ReservationRepository) GetRecentGuests(ctx context.Context) ([]models.GuestInfo, error) {
	oneMonthAgo := time.Now().AddDate(0, -1, 0)
	query := `
        SELECT u.id, u.username, ro.number, res.check_in, res.check_out
        FROM reservations res
        JOIN users u ON res.user_id = u.id
        JOIN rooms ro ON res.room_id = ro.id
        WHERE res.check_in >= $1
        AND res.status IN ('checked_in', 'checked_out')
        ORDER BY res.check_in DESC
    `
	rows, err := r.db.Query(ctx, query, oneMonthAgo)
	if err != nil {
		r.logger.Error("Failed to get recent guests", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var guests []models.GuestInfo
	for rows.Next() {
		var guest models.GuestInfo
		err := rows.Scan(&guest.UserID, &guest.Username, &guest.RoomNumber, &guest.CheckIn, &guest.CheckOut)
		if err != nil {
			r.logger.Error("Failed to scan guest information", zap.Error(err))
			return nil, err
		}
		guests = append(guests, guest)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error scanning guest rows", zap.Error(err))
		return nil, err
	}

	return guests, nil
}
