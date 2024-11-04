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

type RoomRepositoryInterface interface {
	GetAvailableRooms(ctx context.Context, checkIn, checkOut time.Time) ([]models.Room, error)
	GetRoomByID(ctx context.Context, roomID int) (*models.Room, error)
	UpdateRoomAvailability(ctx context.Context, roomID int, isAvailable bool) error
}

type RoomRepository struct {
	db     *pgxpool.Pool
	logger *logger.ZapLogger
}

func NewRoomRepository(db *pgxpool.Pool, logger *logger.ZapLogger) *RoomRepository {
	return &RoomRepository{
		db:     db,
		logger: logger,
	}
}

func (r *RoomRepository) GetAvailableRooms(ctx context.Context, checkIn, checkOut time.Time) ([]models.Room, error) {
	query := `
        SELECT id, hotel_id, number, type, price, is_available, created_at
        FROM rooms
        WHERE is_available = TRUE
        AND id NOT IN (
            SELECT room_id FROM reservations
            WHERE status IN ('reserved', 'checked_in')
            AND ($1 < check_out AND $2 > check_in)
        )
    `
	rows, err := r.db.Query(ctx, query, checkIn, checkOut)
	if err != nil {
		r.logger.Error("Failed to get available rooms", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.HotelID, &room.Number, &room.Type, &room.Price, &room.IsAvailable, &room.CreatedAt)
		if err != nil {
			r.logger.Error("Failed to scan room", zap.Error(err))
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Row error when getting available rooms", zap.Error(err))
		return nil, err
	}

	return rooms, nil
}

func (r *RoomRepository) GetRoomByID(ctx context.Context, roomID int) (*models.Room, error) {
	query := `
        SELECT id, hotel_id, number, type, price, is_available, created_at
        FROM rooms
        WHERE id = $1
    `
	room := &models.Room{}
	err := r.db.QueryRow(ctx, query, roomID).Scan(
		&room.ID,
		&room.HotelID,
		&room.Number,
		&room.Type,
		&room.Price,
		&room.IsAvailable,
		&room.CreatedAt,
	)
	if err != nil {
		r.logger.Error("Failed to get room by ID", zap.Error(err))
		return nil, err
	}
	return room, nil
}

func (r *RoomRepository) UpdateRoomAvailability(ctx context.Context, roomID int, isAvailable bool) error {
	query := `
        UPDATE rooms
        SET is_available = $1
        WHERE id = $2
    `
	cmd, err := r.db.Exec(ctx, query, isAvailable, roomID)
	if err != nil {
		r.logger.Error("Failed to update room availability", zap.Error(err))
		return err
	}
	if cmd.RowsAffected() == 0 {
		errMsg := "room not found"
		r.logger.Warn(errMsg, zap.Int("room_id", roomID))
		return errors.New(errMsg)
	}
	r.logger.Info("Room availability updated", zap.Int("room_id", roomID), zap.Bool("is_available", isAvailable))
	return nil
}

func (r *RoomRepository) GetRoomsByHotel(ctx context.Context, hotelID int) ([]models.Room, error) {
	query := `
        SELECT id, hotel_id, number, type, price, is_available, created_at
        FROM rooms
        WHERE hotel_id = $1
        ORDER BY number ASC
    `
	rows, err := r.db.Query(ctx, query, hotelID)
	if err != nil {
		r.logger.Error("Failed to get rooms by hotel", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.HotelID,
			&room.Number,
			&room.Type,
			&room.Price,
			&room.IsAvailable,
			&room.CreatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan room", zap.Error(err))
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Row error when getting rooms by hotel", zap.Error(err))
		return nil, err
	}

	return rooms, nil
}

func (r *RoomRepository) CreateRoom(ctx context.Context, room *models.Room) error {
	query := `
        INSERT INTO rooms (hotel_id, number, type, price, is_available, created_at)
        VALUES ($1, $2, $3, $4, $5, NOW())
        RETURNING id, created_at
    `
	err := r.db.QueryRow(ctx, query, room.HotelID, room.Number, room.Type, room.Price, room.IsAvailable).Scan(&room.ID, &room.CreatedAt)
	if err != nil {
		r.logger.Error("Failed to create room", zap.Error(err))
		return err
	}
	r.logger.Info("Room created successfully", zap.Int("room_id", room.ID))
	return nil
}
