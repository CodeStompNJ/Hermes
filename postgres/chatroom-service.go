package postgres

import "fmt"

// CreateRoom - create a chatroom
func CreateRoom(name string, description string) {
	sqlStatement := `
	INSERT INTO chatroom (name, description)
	VALUES ($1, $2)
	RETURNING id`
	id := 0
	err := database.QueryRow(sqlStatement, name, description).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)

	return
}

// func GetDefaultRoom() {
// 	sqlStatement := `

// 	`
// }
