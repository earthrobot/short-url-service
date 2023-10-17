package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

// checkResponseCode is a simple utility to check the response code
// of the response
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestCreateShortLinkHandler(t *testing.T) {
	link := "https://ya.ru"

	s := CreateNewServer()

	s.MountHandlers()

	req, _ := http.NewRequest("POST", "/", strings.NewReader(link))

	response := executeRequest(req, s)

	checkResponseCode(t, http.StatusCreated, response.Code)

	require.Equal(t, "text/plain; charset=utf-8", response.Header().Get("Content-Type"))

	require.NotEqual(t, 0, response.Body.Len())
}

func TestGetOriginalLinkHandler(t *testing.T) {
	link := "https://ya.ru"

	s := CreateNewServer()

	s.MountHandlers()

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
