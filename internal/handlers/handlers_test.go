package handlers

import (
	"fmt"
	"github.com/earthrobot/short-url-service/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockStorage struct {
	data map[string]string
}

func (m *mockStorage) Set(key, value string) {
	m.data[key] = value
}

func (m *mockStorage) Get(key string) (string, bool) {
	value, exists := m.data[key]
	return value, exists
}

func executeRequest(req *http.Request, mux *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

//
//func checkResponseCode(t *testing.T, expected, actual int) {
//	if expected != actual {
//		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
//	}
//}

func TestCreateShortLinkHandler(t *testing.T) {
	link := "https://ya.ru"

	//mockDB := &mockStorage{
	//	data: make(map[string]string),
	//}

	store := storage.NewInMemoryStorage()
	handler := NewHandler(store)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/", handler.CreateShortLinkHandler)
	router.Get("/{linkHash}", handler.GetOriginalLinkHandler)

	req, _ := http.NewRequest("POST", "/", strings.NewReader(link))
	response := executeRequest(req, router)

	if http.StatusCreated != response.Code {
		t.Errorf("Expected response code %d. Got %d\n",
			http.StatusCreated, response.Code)
	}

	require.Equal(t, "text/plain; charset=utf-8", response.Header().Get("Content-Type"))
	require.NotEqual(t, 0, response.Body.Len())
}

func TestGetOriginalLinkHandler(t *testing.T) {
	link := "https://ya.ru"

	//mockDB := &mockStorage{
	//	data: make(map[string]string),
	//}

	store := storage.NewInMemoryStorage()
	handler := NewHandler(store)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/", handler.CreateShortLinkHandler)

	reqPost, _ := http.NewRequest("POST", "/", strings.NewReader(link))

	responsePost := executeRequest(reqPost, router)
	shortLink := responsePost.Body.String()

	println("shortLink:", shortLink)
	fullLink := fmt.Sprintf("http://localhost:8080%s", shortLink)
	println("fullLink:", fullLink)

	reqGet, _ := http.NewRequest("GET", fullLink, nil)
	responseGet := executeRequest(reqGet, router)

	if http.StatusTemporaryRedirect != responseGet.Code {
		t.Errorf("Expected response code %d. Got %d\n",
			http.StatusTemporaryRedirect, responseGet.Code)
	}

	location := responseGet.Header().Get("Location")
	if location != link {
		t.Errorf("Get Original Link Handler returned wrong Location header: got %v want %v", location, link)
	}
}
