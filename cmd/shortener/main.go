package main

import (
	"net/http"

	"github.com/earthrobot/short-url-service/config"
	"github.com/earthrobot/short-url-service/internal/router"
)

func main() {

	// инициируем конфиг
	config.Load()

	r := router.NewRouter()

	err := http.ListenAndServe(config.ConfSet.AppHost, r.Router)
	if err != nil {
		panic(err)
	}
}
