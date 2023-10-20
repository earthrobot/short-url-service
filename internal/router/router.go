package router

import (
	"github.com/earthrobot/short-url-service/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	Router *chi.Mux
}

func NewRouter() *Router {
	s := &Router{}
	s.Router = chi.NewRouter()
	db := storage.NewInMemoryStorage()
	h := NewHandler(db)

	s.Router.Use(middleware.Logger)
	s.Router.Post("/", h.createShortLinkHandler)
	s.Router.Get("/{linkHash}", h.getOriginalLinkHandler)

	return s
}
