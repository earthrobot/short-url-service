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
	ConfSet.AppHost = getConfigValue("a", "localhost:8080", "SERVER_ADDRESS", "Host for the app")
	ConfSet.URLHost = getConfigValue("b", "http://localhost:8080", "BASE_URL", "Host for links")
	ConfSet.FileStoragePath = getConfigValue("f", "/tmp/short-url-db.json", "FILE_STORAGE_PATH", "File storage path")
	flag.Parse()
}

func getConfigValue(flagName, defaultValue, envName, description string) string {
	flagValue := flag.String(flagName, defaultValue, description)

	if envValue := os.Getenv(envName); envValue != "" {
		return envValue
	}

	return *flagValue
}
