package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/earthrobot/short-url-service/internal/server"
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

func executeRequest(req *http.Request, s *server.Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestCreateShortLinkHandler(t *testing.T) {
	link := "https://ya.ru"

	mockDB := &mockStorage{
		data: make(map[string]string),
	}

	handler := NewHandler(mockDB)
	s := server.NewServer(handler)

	req, _ := http.NewRequest("POST", "/", strings.NewReader(link))
	response := executeRequest(req, s)

	checkResponseCode(t, http.StatusCreated, response.Code)
	require.Equal(t, "text/plain; charset=utf-8", response.Header().Get("Content-Type"))
	require.NotEqual(t, 0, response.Body.Len())
}

func TestGetOriginalLinkHandler(t *testing.T) {
	link := "https://ya.ru"

	mockDB := &mockStorage{
		data: make(map[string]string),
	}
	s := server.NewServer(mockDB)

	reqPost, _ := http.NewRequest("POST", "/", strings.NewReader(link))
	responsePost := executeRequest(reqPost, s)
	shortLink := responsePost.Body.String()

	reqGet, _ := http.NewRequest("GET", shortLink, nil)
	responseGet := executeRequest(reqGet, s)

	checkResponseCode(t, http.StatusTemporaryRedirect, responseGet.Code)

	location := responseGet.Header().Get("Location")
	if location != link {
		t.Errorf("Get Original Link Handler returned wrong Location header: got %v want %v", location, link)
	}
}
