package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"lr2/internal/auth"
	"net/http"
	"strconv"
	"time"

	"lr2/internal/models"
	"lr2/internal/service"
)

type Handlers struct {
	Service *service.HotelsService
	Store   *sessions.CookieStore
}

func (h *Handlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.Service.RegisterUser(r.Context(), req.Username, req.Email, hashedPassword); err != nil {
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func (h *Handlers) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.Service.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *Handlers) GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email parameter is required", http.StatusBadRequest)
		return
	}

	user, err := h.Service.GetUserByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *Handlers) CreateHotelHandler(w http.ResponseWriter, r *http.Request) {
	var hotel models.Hotel

	if err := json.NewDecoder(r.Body).Decode(&hotel); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateHotel(r.Context(), hotel.Name, hotel.Address, hotel.City, hotel.Country, hotel.Phone, hotel.Email); err != nil {
		http.Error(w, "Failed to create hotel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Hotel created successfully"))
}

func (h *Handlers) GetHotelByIDHandler(w http.ResponseWriter, r *http.Request) {
	hotelIDStr := chi.URLParam(r, "hotelID")
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	hotel, err := h.Service.GetHotelByID(r.Context(), hotelID)
	if err != nil {
		http.Error(w, "Hotel not found: "+err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(hotel)
}

func (h *Handlers) GetAllHotelsHandler(w http.ResponseWriter, r *http.Request) {
	hotels, err := h.Service.GetAllHotels(r.Context())
	if err != nil {
		http.Error(w, "Failed to get hotels: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(hotels)
}

func (h *Handlers) UpdateHotelHandler(w http.ResponseWriter, r *http.Request) {
	hotelIDStr := chi.URLParam(r, "hotelID")
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	var hotel models.Hotel

	if err := json.NewDecoder(r.Body).Decode(&hotel); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hotel.ID = hotelID

	if err := h.Service.UpdateHotel(r.Context(), &hotel); err != nil {
		http.Error(w, "Failed to update hotel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Hotel updated successfully"))
}

func (h *Handlers) DeleteHotelHandler(w http.ResponseWriter, r *http.Request) {
	hotelIDStr := chi.URLParam(r, "hotelID")
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteHotel(r.Context(), hotelID); err != nil {
		http.Error(w, "Failed to delete hotel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Hotel deleted successfully"))
}

func (h *Handlers) GetRoomsByHotelHandler(w http.ResponseWriter, r *http.Request) {
	hotelIDStr := chi.URLParam(r, "hotelID")
	hotelID, err := strconv.Atoi(hotelIDStr)
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	rooms, err := h.Service.GetRoomByHotel(r.Context(), hotelID)
	if err != nil {
		http.Error(w, "Failed to get rooms: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rooms)
}

func (h *Handlers) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		HotelID int     `json:"hotel_id"`
		Number  string  `json:"number"`
		Type    string  `json:"type"`
		Price   float64 `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateRoom(r.Context(), req.HotelID, req.Number, req.Type, req.Price); err != nil {
		http.Error(w, "Failed to create room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Room created successfully"))
}

func (h *Handlers) ReserveRoomHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   int    `json:"user_id"`
		RoomID   int    `json:"room_id"`
		CheckIn  string `json:"check_in"`
		CheckOut string `json:"check_out"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	checkIn, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		http.Error(w, "Invalid check-in date format", http.StatusBadRequest)
		return
	}
	checkOut, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		http.Error(w, "Invalid check-out date format", http.StatusBadRequest)
		return
	}

	if err := h.Service.ReserveRoom(r.Context(), req.UserID, req.RoomID, checkIn, checkOut); err != nil {
		http.Error(w, "Failed to reserve room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Room reserved successfully"))
}

func (h *Handlers) CheckInUserHandler(w http.ResponseWriter, r *http.Request) {
	reservationIDStr := chi.URLParam(r, "reservationID")
	reservationID, err := strconv.Atoi(reservationIDStr)
	if err != nil {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.CheckInUser(r.Context(), reservationID); err != nil {
		http.Error(w, "Failed to check in user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User checked in successfully"))
}

func (h *Handlers) CheckOutUserHandler(w http.ResponseWriter, r *http.Request) {
	reservationIDStr := chi.URLParam(r, "reservationID")
	reservationID, err := strconv.Atoi(reservationIDStr)
	if err != nil {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.CheckOutUser(r.Context(), reservationID); err != nil {
		http.Error(w, "Failed to check out user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User checked out successfully"))
}

func (h *Handlers) WriteReviewHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ReservationID int    `json:"reservation_id"`
		Rating        int    `json:"rating"`
		Comment       string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Service.WriteReview(r.Context(), req.ReservationID, req.Rating, req.Comment); err != nil {
		http.Error(w, "Failed to write review: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Review submitted successfully"))
}

func (h *Handlers) GetRecentGuestsHandler(w http.ResponseWriter, r *http.Request) {
	guests, err := h.Service.GetRecentGuests(r.Context())
	if err != nil {
		http.Error(w, "Failed to get recent guests: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(guests)
}

func (h *Handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.Service.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	session, err := h.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	session.Values["user_id"] = user.ID

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Login successful"))
}

func (h *Handlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := h.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	session.Values["user_id"] = nil
	session.Options.MaxAge = -1

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Logout successful"))
}
