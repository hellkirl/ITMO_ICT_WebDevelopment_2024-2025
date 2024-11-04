package api

import (
	"github.com/go-chi/chi"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	mux        *chi.Mux
}

func NewServer(httpServer *http.Server) *Server {
	return &Server{
		mux:        chi.NewRouter(),
		httpServer: httpServer,
	}
}

func (s *Server) setUpRoutes() {
	//s.mux.Get("/hotels", )
}
