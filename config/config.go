package config

import (
	"flag"
	"os"
)

type Config struct {
	AppHost string
	URLHost string
}

var ConfSet Config

func Load() {

	// дефолтные значения
	ConfSet.AppHost = "localhost:8080"
	ConfSet.URLHost = "http://localhost:8080"

	// парсим флаги
	appHostFlag := flag.String("a", "localhost:8080", "Host for the app")
	urlHostFlag := flag.String("b", "http://localhost:8080", "Host for links")
	flag.Parse()

	// проставляем из флагов или энвов если есть
	if *appHostFlag != "" {
		ConfSet.AppHost = *appHostFlag
	} else if os.Getenv("SERVER_ADDRESS") != "" {
		ConfSet.AppHost = os.Getenv("SERVER_ADDRESS")
	}

	if *urlHostFlag != "" {
		ConfSet.URLHost = *urlHostFlag
	} else if os.Getenv("BASE_URL") != "" {
		ConfSet.URLHost = os.Getenv("BASE_URL")
	}

}
