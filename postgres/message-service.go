package postgres

import "fmt"

/**
id SERIAL PRIMARY KEY,
		chatroom_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		text TEXT NOT NULL,
		created_at TIMESTAMP with time zone DEFAULT current_timestamp
*/
type Message struct {
	ID         int
	chatroomID int
	userID     int
	text       string
}

// CreateMessage - create a new message
func CreateMessage(text string, chatroomID int, userID int) int {
	sqlStatement := `
	INSERT INTO messages (chatroom_id, user_id, text)
	VALUES ($1, $2, $3)
	RETURNING id`
	var id int
	err := database.QueryRow(sqlStatement, chatroomID, userID, text).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)

	return id
}

// GetMessagesForRoom - get all messages for a room
func GetMessagesForRoom(chatroomID int) {
	sqlStatement := `
	SELECT id, user_id, chatroom_id, text FROM messages WHERE chatroom_id=$1
	`
	rows, err := database.Query(sqlStatement, chatroomID)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var message Message
		err = rows.Scan(&message.ID, &message.userID, &message.chatroomID, &message.text)
		if err != nil {
			// handle this error
			panic(err)
		}
		fmt.Println(message)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return
}

// GetMessagesForUser - get all messages for user
func GetMessagesForUser(userID int) {
	sqlStatement := `
	SELECT id, user_id, chatroom_id, text FROM messages WHERE user_id=$1
	`
	rows, err := database.Query(sqlStatement, userID)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var message Message
		err = rows.Scan(&message.ID, &message.userID, &message.chatroomID, &message.text)
		if err != nil {
			// handle this error
			panic(err)
		}
		fmt.Println(message)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return
}

// GetMessagesForUserAndRoom - get all messages for user and room
func GetMessagesForUserAndRoom(userID int, chatroomID int) {
	sqlStatement := `
	SELECT id, user_id, chatroom_id, text FROM messages WHERE user_id=$1 AND chatroom_id=$2
	`
	rows, err := database.Query(sqlStatement, userID, chatroomID)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var message Message
		err = rows.Scan(&message.ID, &message.userID, &message.chatroomID, &message.text)
		if err != nil {
			// handle this error
			panic(err)
		}
		fmt.Println(message)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return
}


//deleting message verification to be added
func DeleteMessage(messageID int) {
	sqlStatement := `
		DELETE FROM messages
		WHERE id = $1;
	`
	rows, err := database.Query(sqlStatement, messageID)
	if err != nil {
		panic(err)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}


}