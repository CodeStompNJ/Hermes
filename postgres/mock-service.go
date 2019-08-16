package postgres

import (
	"fmt"
)

// SetupMockMessages - Used for setting up test messages
func SetupMockMessages() {
	if initialized {

		// room 1
		CreateMessage("test message", 1, 1)
		CreateMessage("test message 2", 1, 1)
		CreateMessage("hello world!", 1, 2)
		CreateMessage("hello world!", 1, 3)

		// room 2
		CreateMessage("test message", 2, 1)
		CreateMessage("test message 2", 2, 1)
		CreateMessage("hello world!", 2, 2)
		CreateMessage("hello world!", 2, 3)

		fmt.Println("\n\nget messages for room 1:")
		GetMessagesForRoom(1)
		fmt.Println("\nget messages for room 2:")
		GetMessagesForRoom(2)

		fmt.Println("\n\nget messages for user 1:")
		GetMessagesForUser(1)
		fmt.Println("\nget messages for user 2:")
		GetMessagesForUser(2)
		fmt.Println("\nget messages for user 3:")
		GetMessagesForUser(3)
		fmt.Println("\nGet message with ID 1:")
		ReturnMessage(1)
	} else {
		fmt.Println("Table can't be initialized: DB is not running!")
	}
}

// CreateMockUsers - mock users
func CreateMockUsers() {
	if initialized {
		fmt.Println("creating users")
		CreateUser("jbond", "James", "Bond", "james.bond@mail.me.the.stuff.com", "password1")
		CreateUser("freeman", "Mr", "Freeman", "freedman@themail-house.com", "password2")
		CreateUser("dennis", "the", "menace", "dman@housing.com", "password3")
		fmt.Println("creating users done")
	} else {
		fmt.Println("Table can't be initialized: DB is not running!")
	}
}
