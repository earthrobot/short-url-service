package main

import (
	"net/http"

	"github.com/earthrobot/short-url-service/config"
	"github.com/earthrobot/short-url-service/internal/router"
	"github.com/earthrobot/short-url-service/internal/storage"
)

func main() {

	config.Load()
	db, _ := storage.NewInMemoryStorage(config.ConfSet.FileStoragePath)
	r := router.NewRouter(db)

	err := http.ListenAndServe(config.ConfSet.AppHost, r.Router)
	if err != nil {
		panic(err)
	}
}
