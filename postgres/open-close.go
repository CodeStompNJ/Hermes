package postgres

import (
	"database/sql"
	"fmt"

	util "../utilities"

	_ "github.com/lib/pq" // db connection
)

var initialized = false
var database *sql.DB

// OpenDBConnection - open connection to db
func OpenDBConnection() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		util.GetEnv("POSTGRES_HOST", "localhost"),
		util.GetEnv("POSTGRES_PORT", "5433"),
		util.GetEnv("HERMES_USER", "seshat"),
		util.GetEnv("HERMES_PASS", "*"),
		util.GetEnv("HERMES_DB_NAME", "hermes"))

	fmt.Println(psqlInfo)

	db, dbErr := sql.Open("postgres", psqlInfo)
	if dbErr != nil {
		panic(dbErr)
	}

	dbErr = db.Ping()
	if dbErr != nil {
		panic(dbErr)
	}

	database = db

	initialized = true
}

// CloseDBConnection - close connection to db
func CloseDBConnection() {
	fmt.Println("Closing Database Connections")

	if initialized {
		database.Close()
		initialized = false
	}

}
