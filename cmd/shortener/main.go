package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type DataStorage struct {
	data map[string]string
}

func NewDataStorage() *DataStorage {
	return &DataStorage{
		data: make(map[string]string),
	}
}

func (ds *DataStorage) Set(key, value string) {
	ds.data[key] = value
}

func (ds *DataStorage) Get(key string) (string, bool) {
	value, exists := ds.data[key]
	return value, exists
}

var dataStorage = NewDataStorage()

func mainHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "Unsupported Media Type", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		link := string(body)

		hash := sha256.Sum256([]byte(link))
		shortLink := fmt.Sprintf("%x", hash)

		dataStorage.Set(shortLink, link)
		w.WriteHeader(http.StatusCreated)
		returnLink := "http://localhost:8080/" + shortLink
		w.Write([]byte(returnLink))

	} else if r.Method == http.MethodGet {

		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 2 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		shortLink := parts[1]
		link, ok := dataStorage.Get(shortLink)
		if !ok {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", link)
		w.WriteHeader(307)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(`/`, mainHandle)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
