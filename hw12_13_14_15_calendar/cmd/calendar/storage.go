package main

import (
	"database/sql"
	"fmt"
	"github.com/IvanSkripnikov/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgx/stdlib"
)

var (
	db   sql.DB
	logg logger.Logger
)

type (
	Event struct {
		ID           int
		Title        string
		DateStart    string
		DateEnd      string
		Description  string
		UserID       int
		RememberTime int
	}

	Notice struct {
		ID     int
		Title  string
		Date   string
		UserID int
	}
)

const migrationsDir = "./migrations"

type DataBaseConnection struct {
	host     string
	db       string
	user     string
	password string
}

// NewDBConnection подключение к БД сервиса
func NewDBConnection() *sql.DB {
	dbConn := DataBaseConnection{
		host:     os.Getenv("DB_HOST"),
		db:       os.Getenv("DB_NAME"),
		user:     os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASS"),
	}
	dataSource := dbConn.getDataSource()
	db, err := sql.Open("mysql", dataSource)
	errPing := db.Ping()

	if err != nil {
		fatalMessage := fmt.Sprintf("Failed to connect to service database. Error: %v", err)
		logg.Fatal(fatalMessage)
	} else if errPing != nil {
		fatalMessage := fmt.Sprintf("Failed to ping service database. Error: %v", errPing)
		logg.Fatal(fatalMessage)
	} else {
		logg.Debug("Connection to the service database successfully established.")
	}

	return db
}

// getDataSource Получить строку соединения для БД
func (conn *DataBaseConnection) getDataSource() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		conn.user, conn.password, conn.host, 3306, conn.db)
}

func main() {
	dsn := "postgres://myuser:mypass@localhost:5432/mydb?sslmode=verify-full"
	db, err := sql.Open("pgx", dsn)
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

// CreateTables Выполнить запросы на создание таблиц
func CreateTables() {
	dbConn := NewDBConnection()
	defer dbConn.Close()

	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		errMessage := fmt.Sprintf("Failed to get list of migration files. Error: %v", err)
		logg.Error(errMessage)
	} else {
		logg.Debug("List of migration files retrieved successfully.")
	}

	for _, file := range files {
		if !file.IsDir() {
			migration := models.Migration{
				Version: file.Name(),
			}

			if !migration.HasExistsRow(dbConn) {
				data, err := ioutil.ReadFile(migrationsDir + "/" + file.Name())

				if err != nil {
					errMessage := fmt.Sprintf("Failed to read migration file: %v. Error: %v", file.Name(), err)
					logg.Error(errMessage)
				} else {
					debugMessage := fmt.Sprintf("The migration file was successfully read: %v.", file.Name())
					logg.Debug(debugMessage)
				}

				sqlQuery := strings.ReplaceAll(string(data), "\r\n", "")
				result, err := dbConn.Exec(sqlQuery)
				migration.InsertRow(dbConn)

				if err == nil && result != nil {
					infoMessage := fmt.Sprintf("Migration has been applied successfully: %v.", file.Name())
					logg.Info(infoMessage)
				}
			}
		}
	}
}
