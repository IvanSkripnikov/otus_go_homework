package controllers

import (
	"app/database"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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
