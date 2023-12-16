package controllers

import (
	"net/http"

	"app/helpers"
)

func BannersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		helpers.GetAllBanners(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func BannerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		helpers.GetBanner(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func AddBannerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		helpers.AddBannerToSlot(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func RemoveBannerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		helpers.RemoveBannerFromSlot(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func GetBannerForShowHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		helpers.GetBannerForShow(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ClickHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		helpers.EventClick(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
