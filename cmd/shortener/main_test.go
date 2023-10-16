package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateShortLinkHandler(t *testing.T) {
	link := "https://ya.ru"
	request, err := http.NewRequest("POST", "/", strings.NewReader(link))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	createShortLinkHandler(recorder, request)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("Create Short Link Handler returned wrong status code: got %v expected %v", status, http.StatusCreated)
	}

	contentType := recorder.Header().Get("Content-Type")
	if contentType != "text/plain; charset=utf-8" {
		t.Errorf("Create Short Link Handler returned wrong Content-Type header: got %v want %v", contentType, "text/plain; charset=utf-8")
	}

	if recorder.Body.Len() == 0 {
		t.Errorf("Create Short Link Handler returned an empty response body, but expected some content")
	}
}

func TestGetOriginalLinkHandler(t *testing.T) {
	link := "https://ya.ru"
	requestPost, err := http.NewRequest("POST", "/", strings.NewReader(link))
	if err != nil {
		t.Fatal(err)
	}
	recorderPost := httptest.NewRecorder()
	createShortLinkHandler(recorderPost, requestPost)

	shortLink := recorderPost.Body.String()
	requestGet, err := http.NewRequest("GET", shortLink, nil)
	if err != nil {
		t.Fatal(err)
	}
	recorderGet := httptest.NewRecorder()
	getOriginalLinkHandler(recorderGet, requestGet)

	if status := recorderGet.Code; status != 307 {
		t.Errorf("Get Original Link Handler returned wrong status code: got %v expected %v", status, 307)
	}

	location := recorderGet.Header().Get("Location")
	if location != link {
		t.Errorf("Get Original Link Handler returned wrong Location header: got %v want %v", location, link)
	}

}
