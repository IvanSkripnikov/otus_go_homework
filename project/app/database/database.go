package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var DB *sql.DB

func init() {
	DB = InitDataBase()
}

func InitDataBase() *sql.DB {
	fmt.Println("connecting ...")

	// get environment variables
	env := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}

	user := env("MYSQL_USER", "user")
	pass := env("MYSQL_PASSWORD", "pass")
	prot := env("MYSQL_PROT", "tcp")
	addr := env("MYSQL_ADDR", "localhost:3306")
	dbname := env("MYSQL_DATABASE", "test")
	netAddr := fmt.Sprintf("%s(%s)", prot, addr)
	dsn := fmt.Sprintf("%s:%s@%s/%s?timeout=30s", user, pass, netAddr, dbname)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("DB connection has been failed.", err.Error())
	}

	fmt.Println("connected!!")

	return db
}
