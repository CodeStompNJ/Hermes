package postgres

import "fmt"

//id
//name
//description
//createdat

// CreateRoom - create a chatroom
func CreateRoom(name string, description string) {
	sqlStatement := `
	INSERT INTO chatroom (name, description)
	VALUES ($1, $2)
	RETURNING id`
	var id int
	err := database.QueryRow(sqlStatement, name, description).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New chatroom ID is:", id)

}

//get latest room
func LatestRoom(userID int) {
	sqlStatement := `
	SELECT c.id, c.name, m.id AS message_id, m.text AS message_text, m.created_at
	FROM chatroom AS c
	INNER JOIN messages AS m ON
	m.chatroom_id=c.id
	ORDER BY m.created_at DESC
	LIMIT 1;
	`
	var id int
	err := database.QueryRow(sqlStatement, userID).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New chatroom ID is:", id)

}

//delete chatroom ID
func DeleteChatroom(idEdit int){
	sqlStatement := `
	DELETE FROM chatroom
	WHERE id = $1;`
	_, err := database.Exec(sqlStatement, idEdit)
	if err != nil {
	  panic(err)
	}

}

//add user to junction table, 
func CreateJoin(userID int, chatroomID int) {
	sqlStatement := `
	INSERT INTO junctionUC (userID, chatroomID)
	VALUES ($1, $2)
	RETURNING id`
	id := 0
	err := database.QueryRow(sqlStatement, userID, chatroomID)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)

}

//delete role from junction table
func DeleteJoin(idEdit int) {
	sqlStatement := `
	DELETE FROM junctionUC
	WHERE is = $1`
	_, err := database.Exec(sqlStatement, idEdit)
	if err != nil {
	  panic(err)
	}

}//no reason to edit ittems in the junction table right now

//edit chatroom

//does chatroom exist
func DoesUserBelongtoGroup(idExist int) bool{
	flag := true

	sqlStatement := `
	SELECT * FROM junctionUC
	WHERE id = $1;`
	res, err := database.Exec(sqlStatement, idEdit)
	if err != nil {
	  panic(err)
	}

	count, err := res.RowsAffected()
	if err != nil {
 		panic(err)
	}
	fmt.Println(count)


	return flag
}

