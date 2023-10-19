package main

import (
	"net/http"

	"github.com/earthrobot/short-url-service/config"
	"github.com/earthrobot/short-url-service/internal/handlers"
	"github.com/earthrobot/short-url-service/internal/server"
	"github.com/earthrobot/short-url-service/internal/storage"
)

func main() {
	config.Load()
	storage := storage.NewInMemoryStorage()
	handler := handlers.NewHandler(storage)
	s := server.NewServer(handler)
	err := http.ListenAndServe(config.ConfSet.AppHost, s.Router)
	if err != nil {
		panic(err)
	}
}
