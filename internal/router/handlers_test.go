package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/earthrobot/short-url-service/internal/models"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/stretchr/testify/require"
)

func executeRequest(req *http.Request, rtr *Router) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	rtr.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestCreateShortLinkHandler(t *testing.T) {
	link := "https://ya.ru"

	rtr := NewRouter()

	req, _ := http.NewRequest("POST", "/", strings.NewReader(link))

	response := executeRequest(req, rtr)

	checkResponseCode(t, http.StatusCreated, response.Code)

	require.Equal(t, "text/plain; charset=utf-8", response.Header().Get("Content-Type"))

	require.NotEqual(t, 0, response.Body.Len())
}

func TestGetOriginalLinkHandler(t *testing.T) {
	link := "https://ya.ru"

	rtr := NewRouter()

	reqPost, _ := http.NewRequest("POST", "/", strings.NewReader(link))

	responsePost := executeRequest(reqPost, rtr)
	shortLink := responsePost.Body.String()

	reqGet, _ := http.NewRequest("GET", shortLink, nil)
	responseGet := executeRequest(reqGet, rtr)

	checkResponseCode(t, http.StatusTemporaryRedirect, responseGet.Code)

	location := responseGet.Header().Get("Location")
	if location != link {
		t.Errorf("Get Original Link Handler returned wrong Location header: got %v want %v", location, link)
	}
}

func TestCreateShortLinkApiHandler(t *testing.T) {
	link := "https://ya.ru"

	rtr := NewRouter()

	shortenURLReq, _ := ffjson.Marshal(&models.ShortenRequest{URL: link})

	req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader(string(shortenURLReq)))

	response := executeRequest(req, rtr)

	checkResponseCode(t, http.StatusCreated, response.Code)

	require.Equal(t, "application/json", response.Header().Get("Content-Type"))

	require.NotEqual(t, 0, response.Body.Len())
}
