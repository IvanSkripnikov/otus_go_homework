package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/controllers"

	"github.com/gavv/httpexpect/v2"
)

func TestRoot(t *testing.T) {
	expected := "{\"message\": \"Hello dear friend! Welcome!\"}"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	controllers.HelloPageHandler(w, req)
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

func TestBanner(t *testing.T) {
	handler := GetHttpHandler()

	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/banners/1").
		Expect().
		Status(http.StatusOK).JSON().IsObject()
}

func TestBanners(t *testing.T) {
	handler := GetHttpHandler()

	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/banners").
		Expect().
		Status(http.StatusOK).JSON().Array().NotEmpty()
}

func TestGetBannerForShow(t *testing.T) {
	handler := GetHttpHandler()

	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/get_banner_for_show/slot=1&group=1").
		Expect().
		Status(http.StatusOK).Raw()
}
