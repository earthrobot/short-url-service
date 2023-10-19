package server

import (
	"github.com/earthrobot/short-url-service/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Router *chi.Mux
}

func NewServer(h *handlers.Handlers) *Server {
	s := &Server{}
	s.Router = chi.NewRouter()

	s.Router.Use(middleware.Logger)
	s.Router.Post("/", h.CreateShortLinkHandler)
	s.Router.Get("/{linkHash}", h.GetOriginalLinkHandler)

	return s
}
