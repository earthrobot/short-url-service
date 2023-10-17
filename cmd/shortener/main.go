package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func (s *Server) createShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	link := string(body)

	hash := sha256.Sum256([]byte(link))
	shortLink := fmt.Sprintf("%x", hash)

	s.Db.Set(shortLink, link)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	returnLink := "http://localhost:8080/" + shortLink
	w.Write([]byte(returnLink))
}

func (s *Server) getOriginalLinkHandler(w http.ResponseWriter, r *http.Request) {
	linkHash := chi.URLParam(r, "linkHash")

	link, ok := s.Db.Get(linkHash)
	if !ok {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	s := CreateNewServer()

	s.MountHandlers()

	err := http.ListenAndServe(":8080", s.Router)
	if err != nil {
		panic(err)
	}
}

type Server struct {
	Router *chi.Mux
	Db     *DataStorage
}

func CreateNewServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	s.Db = NewDataStorage()
	return s
}

func (s *Server) MountHandlers() {
	s.Router.Use(middleware.Logger)
	s.Router.Post("/", s.createShortLinkHandler)
	s.Router.Get("/{linkHash}", s.getOriginalLinkHandler)
}
