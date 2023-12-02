package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"app/database"
)

const tableName string = "tasks"

type Task struct {
	Id        int    `json:"id,omitempty"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var task Task
	task.Id, _ = GetIdFromRequestString(r.URL.Path)

	if task.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "{ \"message\": \"Invalid request GetHandler\"}"+r.URL.Path)
		return
	}

	stmt, err := database.Db.Prepare(fmt.Sprintf("SELECT * from %s WHERE id = ?", tableName))

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	defer stmt.Close()

	if err = stmt.QueryRow(task.Id).Scan(&task.Id, &task.Title, &task.Body, &task.CreatedAt, &task.UpdatedAt); err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "{ \"message\": \"Not Found\"}")
		return
	}

	var buf bytes.Buffer
	je := json.NewEncoder(&buf)

	if err = je.Encode(&task); err != nil {
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

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := database.Db.Query(fmt.Sprintf("SELECT * from %s", tableName))

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

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var task Task
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &task)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if task.Title == "" || task.Body == "" {
		log.Println("invalid request format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	stmt, err := database.Db.Prepare(fmt.Sprintf("INSERT INTO %s (title, body, created_at, updated_at) VALUES (?, ?, ?, ?)", tableName))

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	now := time.Now().Format("2006-01-02 15:04:05")
	task.CreatedAt = now
	task.UpdatedAt = now
	res, err := stmt.Exec(task.Title, task.Body, task.CreatedAt, task.UpdatedAt)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	task.Id = int(id)

	var buf bytes.Buffer
	je := json.NewEncoder(&buf)

	if err = je.Encode(&task); err != nil {
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

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var t Task
	t.Id, _ = GetIdFromRequestString(r.URL.Path)
	if t.Id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "{ \"message\": \"Invalid request\"}")
		return
	}

	var buf bytes.Buffer
	je := json.NewEncoder(&buf)

	if err := je.Encode(&t); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	t.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	if t.Title != "" && t.Body != "" {
		stmt, err := database.Db.Prepare(fmt.Sprintf("UPDATE %s SET title = ?, body = ?, updated_at = ? WHERE id = ?", tableName))
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		defer stmt.Close()

		_, err = stmt.Exec(t.Title, t.Body, t.UpdatedAt, t.Id)

		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
	} else if t.Title != "" {
		stmt, err := database.Db.Prepare(fmt.Sprintf("UPDATE %s SET title = ?, updated_at = ? WHERE id = ?", tableName))
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		defer stmt.Close()

		_, err = stmt.Exec(t.Title, t.UpdatedAt, t.Id)

		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
	} else if t.Body != "" {
		stmt, err := database.Db.Prepare(fmt.Sprintf("UPDATE %s SET body = ?, updated_at = ? WHERE id = ?", tableName))
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		defer stmt.Close()

		_, err = stmt.Exec(t.Body, t.UpdatedAt, t.Id)

		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
	} else {
		log.Println("else section")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "{ \"message\": \"else section\" }")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, t.Id)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := GetIdFromRequestString(r.URL.Path)
	stmt, err := database.Db.Prepare(fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName))

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, id)
}
