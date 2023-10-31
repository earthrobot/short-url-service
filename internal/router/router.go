package router

import (
	"github.com/earthrobot/short-url-service/internal/encoding"
	"github.com/earthrobot/short-url-service/internal/logger"
	"github.com/earthrobot/short-url-service/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	Router *chi.Mux
}

func NewRouter() *Router {
	s := &Router{}
	s.Router = chi.NewRouter()
	db := storage.NewInMemoryStorage()
	h := NewHandler(db)

	logger.Initialize()
	s.Router.Use(logger.RequestLogger)
	s.Router.Use(encoding.GzipMiddleware)
	s.Router.Post("/", h.createShortLinkHandler)
	s.Router.Get("/{linkHash}", h.getOriginalLinkHandler)
	s.Router.Post("/api/shorten", h.createShortLinkAPIHandler)

	return s
}
