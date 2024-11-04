// models/models.go
package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type Hotel struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	Phone     string    `json:"phone,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Room struct {
	ID          int       `json:"id"`
	HotelID     int       `json:"hotel_id"`
	Number      string    `json:"number"`
	Type        string    `json:"type"`
	Price       float64   `json:"price"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
}

type Reservation struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	RoomID    int       `json:"room_id"`
	CheckIn   time.Time `json:"check_in"`
	CheckOut  time.Time `json:"check_out"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Review struct {
	ID            int       `json:"id"`
	ReservationID int       `json:"reservation_id"`
	Rating        int       `json:"rating"`
	Comment       string    `json:"comment"`
	CreatedAt     time.Time `json:"created_at"`
}

type GuestInfo struct {
	UserID     int       `json:"user_id"`
	Username   string    `json:"username"`
	RoomNumber string    `json:"room_number"`
	CheckIn    time.Time `json:"check_in"`
	CheckOut   time.Time `json:"check_out"`
}
