package router

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/earthrobot/short-url-service/config"
	"github.com/earthrobot/short-url-service/internal/models"
	"github.com/earthrobot/short-url-service/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/pquerna/ffjson/ffjson"
)

type Handlers struct {
	DB *storage.InMemoryStorage
}

func NewHandler(db *storage.InMemoryStorage) *Handlers {
	return &Handlers{DB: db}
}

func (h *Handlers) createShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	link := string(body)
	hash := sha256.Sum256([]byte(link))
	shortLink := fmt.Sprintf("%x", hash)

	h.DB.Set(shortLink, link)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	returnLink := config.ConfSet.URLHost + "/" + shortLink
	w.Write([]byte(returnLink))
}

func (h *Handlers) getOriginalLinkHandler(w http.ResponseWriter, r *http.Request) {
	linkHash := chi.URLParam(r, "linkHash")
	link, ok := h.DB.Get(linkHash)
	if !ok {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handlers) createShortLinkAPIHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	sr := models.ShortenRequest{}
	err = ffjson.Unmarshal(body, &sr)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	link := sr.URL
	hash := sha256.Sum256([]byte(link))
	shortLink := fmt.Sprintf("%x", hash)

	h.DB.Set(shortLink, link)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	returnLink := config.ConfSet.URLHost + "/" + shortLink

	shRes := models.ShortenResponse{Result: returnLink}
	res, err := ffjson.Marshal(&shRes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(res)
}
