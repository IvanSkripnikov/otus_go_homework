package controllers

import (
	"fmt"
	"log"
	"net/http"
)

func HelloPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		_, err := fmt.Fprint(w, "{\"message\": \"Hello dear friend! Welcome!\"}")
		if err != nil {
			log.Println(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
