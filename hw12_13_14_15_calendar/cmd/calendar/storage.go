package main

import (
	"database/sql"
	"log"
)

var db sql.DB

type Event struct {
	ID           int
	Title        string
	DateStart    string
	DateEnd      string
	Description  string
	UserID       int
	RememberTime int
}

func main() {
	db, err := sql.Open("mysql",
		"user:password@tcp(127.0.0.1:3306)/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
}

func getEvent(id int) (Event, error) {
	var event Event

	query := `SELECT * FROM events WHERE id = $1`
	row := db.QueryRowContext(nil, query, id)

	if err := row.Err(); err != nil {
		return event, err
	}

	var idRecord, userId, rememberTime int
	var title, description, dateStart, dateEnd string

	err := row.Scan(&idRecord, &title, &dateStart, &dateEnd, &description, &userId, &rememberTime)

	if err == sql.ErrNoRows {
		return event, err
	} else if err != nil {
		return event, err
	}

	event.ID = idRecord
	event.Title = title
	event.DateStart = dateStart
	event.DateEnd = dateEnd
	event.Description = description
	event.UserID = userId
	event.RememberTime = rememberTime

	return event, nil
}

func getEventsList() ([]Event, error) {
	var event Event
	var events []Event

	query := `SELECT * FROM events WHERE id > 0`
	rows, err := db.QueryContext(nil, query)
	if err != nil {
		return events, err
	}

	if err := rows.Err(); err != nil {
		return events, err
	}

	defer rows.Close()
	for rows.Next() {
		var idRecord, userId, rememberTime int
		var title, description, dateStart, dateEnd string
		if err := rows.Scan(&idRecord, &title, &dateStart, &dateEnd, &description, &userId, &rememberTime); err != nil {
			return events, err
		}

		event.ID = idRecord
		event.Title = title
		event.DateStart = dateStart
		event.DateEnd = dateEnd
		event.Description = description
		event.UserID = userId
		event.RememberTime = rememberTime

		events = append(events, event)
	}

	return events, nil
}

func deleteEvent(id int) bool {
	// создаем подготовленный запрос
	stmt, err := db.Prepare("DELETE FROM events WHERE id = $1") // *sql.Stmt
	if err != nil {
		log.Fatal(err)
	}
	// освобождаем ресурсы в СУБД
	defer stmt.Close()
	// многократно выполняем запрос

	_, err = stmt.Exec(id)
	if err != nil {
		return false
	}

	return true
}

// передавать сюда мдель с изменёнными данными
func updateEvent(id int, event Event) bool {
	// создаем подготовленный запрос
	query := `UPDATE events 
SET id = $1, title = $2, date_start = $3, date_end = $4, description = $5, user_id = $6, remember_time = $7
WHERE id = $1`
	stmt, err := db.Prepare(query) // *sql.Stmt
	if err != nil {
		log.Fatal(err)
	}
	// освобождаем ресурсы в СУБД
	defer stmt.Close()

	var idRecord, userId, rememberTime int
	var title, description, dateStart, dateEnd string

	idRecord = event.ID
	title = event.Title
	dateStart = event.DateStart
	dateEnd = event.DateEnd
	description = event.Description
	userId = event.UserID
	rememberTime = event.RememberTime

	_, err = stmt.Exec(idRecord, title, dateStart, dateEnd, description, userId, rememberTime)
	if err != nil {
		return false
	}

	return true
}
