package postgres

import (
	"fmt"
)

// Setup - Used for setting up base postgres tables
// TODO move out sql statements for init to sql files for easier upkeep
func SetupDB() {
	if initialized {
		createTables()
		CreateRoom("default", "a default room")
		CreateRoom("extra", "an extra room")
		fmt.Println("creating users now")
		CreateMockUsers()
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
	_, err := database.Exec("DROP TABLE IF EXISTS chatroom cascade")
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
	 *Creating User Table
	 **/
	fmt.Println("creating users table...")
	_, err = database.Exec("DROP TABLE IF EXISTS users cascade")
	if err != nil {
		panic(err)
	}

	createQry = `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username TEXT UNIQUE NOT NULL,
        firstname TEXT,
        lastname TEXT,
		email TEXT NOT NULL,
		password TEXT,
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
	_, err = database.Exec("DROP TABLE IF EXISTS messages cascade")
	if err != nil {
		panic(err)
	}

	createQry = `
	 CREATE TABLE IF NOT EXISTS messages (
		 id SERIAL PRIMARY KEY,
		 chatroom_id INTEGER NOT NULL,
		 user_id INTEGER NOT NULL,
		 text TEXT NOT NULL,
		 created_at TIMESTAMP with time zone DEFAULT current_timestamp,
		 FOREIGN KEY (chatroom_id) REFERENCES chatroom (id),
		 FOREIGN KEY (user_id) REFERENCES users (id)
	 )`
	_, err = database.Exec(createQry)
	if err != nil {
		panic(err)
	}

	//junction table for users and chatroom for many to many relationship
	fmt.Println("creating user/chatroom junction table...")
	_, err = database.Exec("DROP TABLE IF EXISTS junctionUC")
	if err != nil {
		panic(err)
	}

	createQry = `
	CREATE TABLE IF NOT EXISTS junctionUC (
		
		    UserID int NOT NULL,
		    ChatroomID int NOT NULL,
		    CONSTRAINT PK_junctionUC PRIMARY KEY
		    (
		        UserID,
		        ChatroomID
		    ),
		    FOREIGN KEY (UserID) REFERENCES users (id),
		    FOREIGN KEY (ChatroomID) REFERENCES chatroom (id)
		)

	`
	_, err = database.Exec(createQry)
	if err != nil {
		panic(err)
	}

	return
}
