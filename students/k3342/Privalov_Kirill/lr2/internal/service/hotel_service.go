package service

import (
	"lr2/internal/models"
	"lr2/internal/repository"
)

import (
	"context"
	"errors"
	"time"
)

type HotelsService struct {
	repo *repository.HotelsRepository
}

func NewHotelsService(repo *repository.HotelsRepository) *HotelsService {
	return &HotelsService{
		repo: repo,
	}
}

func (s *HotelsService) RegisterUser(ctx context.Context, username, email, password string) error {
	user := &models.User{
		Username: username,
		Email:    email,
		Password: password,
	}
	return s.repo.User.RegisterUser(ctx, user)
}

func (s *HotelsService) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	return s.repo.User.GetUserByID(ctx, userID)
}

func (s *HotelsService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.User.GetUserByEmail(ctx, email)
}

func (s *HotelsService) CreateHotel(ctx context.Context, name, address, city, country, phone, email string) error {
	hotel := &models.Hotel{
		Name:    name,
		Address: address,
		City:    city,
		Country: country,
		Phone:   phone,
		Email:   email,
	}
	return s.repo.Hotel.CreateHotel(ctx, hotel)
}

func (s *HotelsService) GetHotelByID(ctx context.Context, hotelID int) (*models.Hotel, error) {
	return s.repo.Hotel.GetHotelByID(ctx, hotelID)
}

func (s *HotelsService) GetAllHotels(ctx context.Context) ([]models.Hotel, error) {
	return s.repo.Hotel.GetAllHotels(ctx)
}

func (s *HotelsService) UpdateHotel(ctx context.Context, hotel *models.Hotel) error {
	return s.repo.Hotel.UpdateHotel(ctx, hotel)
}

func (s *HotelsService) DeleteHotel(ctx context.Context, hotelID int) error {
	return s.repo.Hotel.DeleteHotel(ctx, hotelID)
}

func (s *HotelsService) GetRoomByHotel(ctx context.Context, hotelID int) ([]models.Room, error) {
	return s.repo.Room.GetRoomsByHotel(ctx, hotelID)
}

func (s *HotelsService) CreateRoom(ctx context.Context, hotelID int, number, roomType string, price float64) error {
	room := &models.Room{
		HotelID:     hotelID,
		Number:      number,
		Type:        roomType,
		Price:       price,
		IsAvailable: true,
	}
	return s.repo.Room.CreateRoom(ctx, room)
}

func (s *HotelsService) ReserveRoom(ctx context.Context, userID, roomID int, checkIn, checkOut time.Time) error {
	user, err := s.repo.User.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	reservation := &models.Reservation{
		UserID:   userID,
		RoomID:   roomID,
		CheckIn:  checkIn,
		CheckOut: checkOut,
		Status:   "reserved",
	}
	return s.repo.Reservation.CreateReservation(ctx, reservation)
}

func (s *HotelsService) CheckInUser(ctx context.Context, reservationID int) error {
	return s.repo.Reservation.CheckIn(ctx, reservationID)
}

func (s *HotelsService) CheckOutUser(ctx context.Context, reservationID int) error {
	return s.repo.Reservation.CheckOut(ctx, reservationID)
}

func (s *HotelsService) WriteReview(ctx context.Context, reservationID, rating int, comment string) error {
	reservation, err := s.repo.Reservation.GetReservationByID(ctx, reservationID)
	if err != nil || reservation == nil {
		return errors.New("reservation not found")
	}
	if reservation.Status != "checked_out" {
		return errors.New("cannot write review for an active reservation")
	}
	review := &models.Review{
		ReservationID: reservationID,
		Rating:        rating,
		Comment:       comment,
	}
	return s.repo.Review.CreateReview(ctx, review)
}

func (s *HotelsService) GetRecentGuests(ctx context.Context) ([]models.GuestInfo, error) {
	return s.repo.Reservation.GetRecentGuests(ctx)
}
