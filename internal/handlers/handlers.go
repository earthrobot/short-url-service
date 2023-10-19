package handlers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/earthrobot/short-url-service/config"
	"github.com/earthrobot/short-url-service/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	DB *storage.InMemoryStorage
}

func NewHandler(db *storage.InMemoryStorage) *Handlers {
	return &Handlers{DB: db}
}

func (h *Handlers) CreateShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	link := string(body)
	hash := sha256.Sum256([]byte(link))
	shortLink := fmt.Sprintf("%x", hash)

	h.DB.Set(shortLink, link)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	returnLink := config.ConfSet.URLHost + "/" + shortLink
	w.Write([]byte(returnLink))
}

func (h *Handlers) GetOriginalLinkHandler(w http.ResponseWriter, r *http.Request) {
	linkHash := chi.URLParam(r, "linkHash")
	link, ok := h.DB.Get(linkHash)
	if !ok {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
