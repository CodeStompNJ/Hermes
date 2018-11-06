package postgres

import "fmt"

//Create a user and add their info to the DB
func CreateUser(username string, firstname string, lastname string, email string) {
	sqlStatement := `
	INSERT INTO users (username, firstname, lastname, email)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	var id int
	err := database.QueryRow(sqlStatement, username, firstname, lastname, email).Scan(&id)
	if err != nil {
		fmt.Println("user failed to create: ", sqlStatement)
		panic(err)
	}
	fmt.Println("New user record ID is:", id)

}

func EditUser(idEdit int, usernameEdit string, firstnameEdit string, lastnameEdit string, emailEdit string) {
	sqlStatement := `
	UPDATE users
	SET username = $2, firstname = $3, lastname = $4, email = $5
	WHERE id = $1;`
	id := 0
	_, err := database.Exec(sqlStatement, idEdit, usernameEdit, firstnameEdit, lastnameEdit, emailEdit)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)

}

func DeleteUser(idEdit int) {
	sqlStatement := `
	DELETE FROM users
	WHERE id = $1;`
	_, err := database.Exec(sqlStatement, idEdit)
	if err != nil {
		panic(err)
	}

}

func DoesUserExist(idEdit int) bool {

	flag := true

	sqlStatement := `
	SELECT * FROM users
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
