package postgres

import (
	"fmt"
)

// Setup - Used for setting up base postgres tables
// TODO move out sql statements for init to sql files for easier upkeep
func SetupDB() {
	if initialized {
		createTables()
		CreateRoom("default", "a default table")
		SetupMockMessages()
	} else {
		fmt.Println("Table can't be initialized: DB is not running!")
	}
}

func createTables() {
	/**
	* Chatroom setup
	 */
	fmt.Println("creating chatroom table...")
	_, err := database.Exec("DROP TABLE IF EXISTS chatroom")
	if err != nil {
		panic(err)
	}

	createQry := `
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

	/**
	* Messages setup
	 */
	fmt.Println("creating messages table...")
	_, err = database.Exec("DROP TABLE IF EXISTS messages")
	if err != nil {
		panic(err)
	}

	createQry = `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		chatroom_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		text TEXT NOT NULL,
		created_at TIMESTAMP with time zone DEFAULT current_timestamp
	)`
	_, err = database.Exec(createQry)
	if err != nil {
		panic(err)
	}
	return
}
