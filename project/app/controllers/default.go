package controllers

import (
	"fmt"
	"log"
	"net/http"
)

func HelloPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err := fmt.Fprint(w, "{ \"message\": \"Hello dear friend! Welcome!\"}")
	if err != nil {
		log.Println(err.Error())
		return
	}
}
