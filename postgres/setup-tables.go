package postgres

import (
	"fmt"
)

// Setup - Used for setting up base postgres tables
// TODO move out sql statements for init to sql files for easier upkeep
func SetupDB() {
	if initialized {
		createTable()
		CreateRoom("default", "a default table")
	} else {
		fmt.Println("Table can't be initialized: DB is not running!")
	}
}

func createTable() {
	fmt.Println("creating chatroom table...")
	_, err := database.Exec("DROP TABLE IF EXISTS chatroom")
	if err != nil {
		panic(err)
	}

	const createQry = `
	CREATE TABLE IF NOT EXISTS chatroom (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		created_at TIMESTAMP with time zone DEFAULT current_timestamp
	)`
	_, err = database.Exec(createQry)
	if err != nil {
		panic(err)
	}
	return
}
