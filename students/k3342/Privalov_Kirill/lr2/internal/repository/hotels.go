package repository

import (
	"go.uber.org/zap"
	"lr2/internal/dbcontext"
	"lr2/pkg/logger"
)

type HotelsRepositoryInterface struct {
	Users        UserRepository
	Hotels       HotelRepository
	Rooms        RoomRepository
	Reservations ReservationRepository
	Reviews      ReviewRepository
}

type HotelsRepository struct {
	User        *UserRepository
	Hotel       *HotelRepository
	Room        *RoomRepository
	Reservation *ReservationRepository
	Review      *ReviewRepository
}

func NewHotelsRepository(dbCtx *dbcontext.DbContext, logger *logger.ZapLogger) *HotelsRepository {
	pool, err := dbCtx.Pool()
	if err != nil {
		logger.Error("couldn't init hotels repository", zap.Error(err))
		return nil
	}
	return &HotelsRepository{
		User:        NewUserRepository(pool, logger),
		Hotel:       NewHotelRepository(pool, logger),
		Room:        NewRoomRepository(pool, logger),
		Reservation: NewReservationRepository(pool, logger),
		Review:      NewReviewRepository(pool, logger),
	}
}
