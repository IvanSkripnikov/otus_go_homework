package main

import (
	"app/controllers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoot(t *testing.T) {
	expected := "{\"message\": \"Hello dear friend! Welcome!\"}"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	controllers.HelloPage(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if string(data) != expected {
		t.Errorf("Expected root message but got %v", string(data))
	}
}

func TestBanners(t *testing.T) {
	expected := "{\"message\": \"Hello dear friend! Welcome!\"}"
	req := httptest.NewRequest(http.MethodGet, "/banners", nil)
	w := httptest.NewRecorder()
	controllers.GetAllBanners(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if string(data) != expected {
		t.Errorf("Expected banners message but got %v", string(data))
	}
}
