package config

import (
	"flag"
	"os"
)

type Config struct {
	AppHost         string
	URLHost         string
	FileStoragePath string
}

var ConfSet Config

func Load() {
	fileStoragePath := flag.String("f", "/tmp/short-url-db.json", "File storage path")
	appHost := flag.String("a", "localhost:8080", "Host for the app")
	urlHost := flag.String("b", "http://localhost:8080", "Host for links")

	flag.Parse()

	ConfSet.FileStoragePath = getEnvOrDefault("FILE_STORAGE_PATH", *fileStoragePath)
	ConfSet.AppHost = getEnvOrDefault("SERVER_ADDRESS", *appHost)
	ConfSet.URLHost = getEnvOrDefault("BASE_URL", *urlHost)
}

func getEnvOrDefault(envName, defaultValue string) string {
	if envValue := os.Getenv(envName); envValue != "" {
		return envValue
	}
	return defaultValue
}
