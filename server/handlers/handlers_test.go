package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	var i HandlerFunc
	var m Marker
	var l []Lab
	i = &Handler{"../../templates", m, l}
	handler := http.HandlerFunc(i.HandleIndex)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestNotExistingPath(t *testing.T) {
	req, err := http.NewRequest("GET", "/fakepath", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	var i HandlerFunc
	var m Marker
	var l []Lab
	i = &Handler{"../../templates", m, l}
	handler := http.HandlerFunc(i.HandleIndex)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "404 page not found"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
