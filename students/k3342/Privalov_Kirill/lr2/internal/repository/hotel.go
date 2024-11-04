package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"lr2/internal/models"
	"lr2/pkg/logger"
)

type HotelRepositoryInterface interface {
	CreateHotel(ctx context.Context, hotel *models.Hotel) error
	GetHotelByID(ctx context.Context, hotelID int) (*models.Hotel, error)
	GetAllHotels(ctx context.Context) ([]models.Hotel, error)
	UpdateHotel(ctx context.Context, hotel *models.Hotel) error
	DeleteHotel(ctx context.Context, hotelID int) error
}

type HotelRepository struct {
	db     *pgxpool.Pool
	logger *logger.ZapLogger
}

func NewHotelRepository(db *pgxpool.Pool, logger *logger.ZapLogger) *HotelRepository {
	return &HotelRepository{
		db:     db,
		logger: logger,
	}
}

func (r *HotelRepository) CreateHotel(ctx context.Context, hotel *models.Hotel) error {
	query := `
        INSERT INTO hotels (name, address, city, country, phone, email, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, NOW())
        RETURNING id, created_at
    `
	err := r.db.QueryRow(ctx, query, hotel.Name, hotel.Address, hotel.City, hotel.Country, hotel.Phone, hotel.Email).Scan(&hotel.ID, &hotel.CreatedAt)
	if err != nil {
		r.logger.Error("Failed to create hotel", zap.Error(err))
		return err
	}
	r.logger.Info("Hotel created successfully", zap.Int("hotel_id", hotel.ID))
	return nil
}

func (r *HotelRepository) GetHotelByID(ctx context.Context, hotelID int) (*models.Hotel, error) {
	query := `
        SELECT id, name, address, city, country, phone, email, created_at
        FROM hotels
        WHERE id = $1
    `
	hotel := &models.Hotel{}
	err := r.db.QueryRow(ctx, query, hotelID).Scan(
		&hotel.ID,
		&hotel.Name,
		&hotel.Address,
		&hotel.City,
		&hotel.Country,
		&hotel.Phone,
		&hotel.Email,
		&hotel.CreatedAt,
	)
	if err != nil {
		r.logger.Error("Failed to get hotel by ID", zap.Error(err))
		return nil, err
	}
	return hotel, nil
}

func (r *HotelRepository) GetAllHotels(ctx context.Context) ([]models.Hotel, error) {
	query := `
        SELECT id, name, address, city, country, phone, email, created_at
        FROM hotels
        ORDER BY created_at DESC
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error("Failed to get all hotels", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var hotels []models.Hotel
	for rows.Next() {
		var hotel models.Hotel
		err := rows.Scan(
			&hotel.ID,
			&hotel.Name,
			&hotel.Address,
			&hotel.City,
			&hotel.Country,
			&hotel.Phone,
			&hotel.Email,
			&hotel.CreatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan hotel", zap.Error(err))
			return nil, err
		}
		hotels = append(hotels, hotel)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Row error when getting all hotels", zap.Error(err))
		return nil, err
	}

	return hotels, nil
}

func (r *HotelRepository) UpdateHotel(ctx context.Context, hotel *models.Hotel) error {
	query := `
        UPDATE hotels
        SET name = $1, address = $2, city = $3, country = $4, phone = $5, email = $6
        WHERE id = $7
    `
	cmd, err := r.db.Exec(ctx, query, hotel.Name, hotel.Address, hotel.City, hotel.Country, hotel.Phone, hotel.Email, hotel.ID)
	if err != nil {
		r.logger.Error("Failed to update hotel", zap.Error(err))
		return err
	}
	if cmd.RowsAffected() == 0 {
		errMsg := "hotel not found"
		r.logger.Warn(errMsg, zap.Int("hotel_id", hotel.ID))
		return errors.New(errMsg)
	}
	r.logger.Info("Hotel updated successfully", zap.Int("hotel_id", hotel.ID))
	return nil
}

func (r *HotelRepository) DeleteHotel(ctx context.Context, hotelID int) error {
	query := `
        DELETE FROM hotels
        WHERE id = $1
    `
	cmd, err := r.db.Exec(ctx, query, hotelID)
	if err != nil {
		r.logger.Error("Failed to delete hotel", zap.Error(err))
		return err
	}
	if cmd.RowsAffected() == 0 {
		errMsg := "hotel not found"
		r.logger.Warn(errMsg, zap.Int("hotel_id", hotelID))
		return errors.New(errMsg)
	}
	r.logger.Info("Hotel deleted successfully", zap.Int("hotel_id", hotelID))
	return nil
}
