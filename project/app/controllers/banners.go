package controllers

import (
	"app/database"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Banner struct {
	Id        int
	Title     string
	Body      string
	CreatedAt string
	Active    bool
}

func GetAllBanners(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := database.Db.Query(fmt.Sprintf("SELECT * from %s", "banners"))

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		task := Task{}
		if err = rows.Scan(&task.Id, &task.Title, &task.Body, &task.CreatedAt, &task.UpdatedAt); err != nil {
			log.Println(err.Error())
			continue
		}
		tasks = append(tasks, task)
	}

	var buf bytes.Buffer
	je := json.NewEncoder(&buf)

	if err = je.Encode(&tasks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprint(w, buf.String())
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func GetBanner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var banner Banner
	banner.Id, _ = GetIdFromRequestString(r.URL.Path)

	if banner.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "{ \"message\": \"Invalid request GetHandler\"}"+r.URL.Path)
		return
	}

	stmt, err := database.Db.Prepare(fmt.Sprintf("SELECT * from %s WHERE id = ?", "banners"))

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	defer stmt.Close()

	if err = stmt.QueryRow(banner.Id).Scan(&banner.Id, &banner.Title, &banner.Body, &banner.CreatedAt, &banner.Active); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "{ \"message\": \"Not Found\"}")
		return
	}

	var buf bytes.Buffer
	je := json.NewEncoder(&buf)

	if err = je.Encode(&banner); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprint(w, buf.String())
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func AddBannerToSlot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bannerId, slotId, resultString := GetIdsFromQueryString(r.URL.Path)

	if resultString != "" {
		log.Println(resultString)
		fmt.Fprint(w, resultString)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stmt, err := database.Db.Prepare(fmt.Sprintf("INSERT INTO %s (banner_id, slot_id) VALUES (?, ?)", "relations_banner_slot"))

	if err != nil {
		log.Println(err.Error())
		fmt.Fprint(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(bannerId, slotId)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprint(w, "{ \"message\": \"Successfully added!\"}")
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func GetIdFromRequestString(url string) (int, error) {
	vars := strings.Split(url, "/")

	return strconv.Atoi(vars[len(vars)-1])
}

func GetIdsFromQueryString(url string) (int, int, string) {
	bannerId := 0
	slotId := 0
	outMessage := ""
	queryParams := strings.Split(url, "/")

	params := strings.Split(queryParams[len(queryParams)-1], "&")
	if len(params) == 1 {
		outMessage = "not all params is set"
		return bannerId, slotId, outMessage
	}

	for _, v := range params {
		pair := strings.Split(v, "=")
		if len(pair) == 1 {
			outMessage = "incorrect params value: " + v
			return bannerId, slotId, outMessage
		}
		if pair[0] == "banner" {
			bannerId, _ = strconv.Atoi(pair[1])
		}
		if pair[0] == "slot" {
			slotId, _ = strconv.Atoi(pair[1])
		}
	}

	return bannerId, slotId, outMessage
}
