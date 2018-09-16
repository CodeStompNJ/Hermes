package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "seshat"
	password = "r8*W6F#8xE"
	dbname   = "hermes"
)

var initialized = false
var database *sql.DB

// OpenDBConnection - open connection to db
func OpenDBConnection() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

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
