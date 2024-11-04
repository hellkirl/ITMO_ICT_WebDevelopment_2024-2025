package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"lr2/internal/dbcontext"
	"lr2/internal/handlers"
	"lr2/internal/repository"
	"lr2/internal/service"
	"lr2/pkg/logger"
	"net/http"
	"os"
)

func main() {
	log, err := logger.NewZapLogger()
	if err != nil {
		fmt.Errorf("couldn't init logger")
	}
	defer log.Sync()

	connURL := "postgres://user:password@localhost:5432/hotelsdb?sslmode=disable"
	dbContext, err := dbcontext.NewDbContext(connURL, log)
	if err != nil {
		log.Error("Failed to connect to database", zap.Error(err))
		os.Exit(1)
	}
	defer dbContext.Close()

	hotelsRepo := repository.NewHotelsRepository(dbContext, log)

	hotelsService := service.NewHotelsService(hotelsRepo)

	handler := handlers.Handlers{
		Service: hotelsService,
	}

	store := sessions.NewCookieStore([]byte("super-secret-key"))

	handler.Store = store

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/users/register", handler.RegisterUserHandler)
	r.Post("/users/login", handler.LoginHandler)

	r.Group(func(r chi.Router) {
		r.Get("/users/{userID}", handler.GetUserByIDHandler)
		r.Get("/users", handler.GetUserByEmailHandler)
		r.Post("/hotels", handler.CreateHotelHandler)
		r.Get("/hotels/{hotelID}", handler.GetHotelByIDHandler)
		r.Get("/hotels", handler.GetAllHotelsHandler)
		r.Put("/hotels/{hotelID}", handler.UpdateHotelHandler)
		r.Delete("/hotels/{hotelID}", handler.DeleteHotelHandler)
		r.Get("/hotels/{hotelID}/rooms", handler.GetRoomsByHotelHandler)
		r.Post("/rooms", handler.CreateRoomHandler)
		r.Post("/reservations", handler.ReserveRoomHandler)
		r.Post("/reservations/{reservationID}/checkin", handler.CheckInUserHandler)
		r.Post("/reservations/{reservationID}/checkout", handler.CheckOutUserHandler)
		r.Post("/reviews", handler.WriteReviewHandler)
		r.Get("/guests/recent", handler.GetRecentGuestsHandler)
	})

	log.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Error("Failed to start server", zap.Error(err))
		os.Exit(1)
	}
}
