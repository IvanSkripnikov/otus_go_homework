package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	user   string
	pass   string
	dbName string
	Db     *sql.DB
)

func init() {
	fmt.Println("connecting ...")

	/*
		user = os.Getenv("MYSQL_USER")
		pass = os.Getenv("MYSQL_PASSWORD")
		dbName = os.Getenv("MYSQL_DATABASE")
	*/

	user = "user"
	pass = "pass"
	dbName = "test"

	if user == "" || pass == "" || dbName == "" {
		log.Fatal("user or pass or dbName does not found.")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(db:3306)/%s", user, pass, dbName)

	var err error
	Db, err = sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("DB connection has been failed.", err.Error())
	}

	fmt.Println("connected!!")
}
