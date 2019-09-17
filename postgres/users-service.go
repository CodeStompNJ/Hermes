package postgres

import (
	"fmt"
	"github.com/lib/pq"
)

type User struct {
	ID        int
	Username  string
	Password  string
	Firstname string
	Lastname  string
	Email     string
	Timestamp string
}

// Create a user and add their info to the DB
func CreateUser(username string, firstname string, lastname string, email string, password string) string {
	result := "success"
	sqlStatement := `
	INSERT INTO users (username, firstname, lastname, email, password)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`
	var id int
	err := database.QueryRow(sqlStatement, username, firstname, lastname, email, password).Scan(&id)
	if err != nil {
		fmt.Println("user failed to create: ", sqlStatement)
		//panic(err.Error)
	}
	if err, ok := err.(*pq.Error); ok {
		fmt.Println(" pq error:", err.Code.Name())
		result = err.Code.Name()
	}

	fmt.Println("New user record ID is:", id)

	return result

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

func ValidUser(username string, password string) bool {
	sqlStatement := `
	SELECT * FROM users
	where username = $1 AND password = $2;
	`
	row, err := database.Query(sqlStatement, username, password)
	if err != nil {
		panic(err)
	}

	var s User
	for row.Next() {
		err = row.Scan(&s.ID, &s.Username, &s.Firstname, &s.Lastname, &s.Email, &s.Timestamp, &s.Password)
		if err != nil {
			// handle this error
			panic(err)
		}
	}

	if s.ID != 0 {
		return true
	}

	return false
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

func ReturnUser(idEdit int) User {
	sqlStatement := `
	SELECT *
	FROM users
	WHERE id = $1;`
	row, err := database.Query(sqlStatement, idEdit)
	if err != nil {
		panic(err)
	}

	var s User
	for row.Next() {
		err = row.Scan(&s.ID, &s.Username, &s.Firstname, &s.Lastname, &s.Email, &s.Timestamp)
		if err != nil {
			// handle this error
			panic(err)
		}
	}

	fmt.Println(s)
	return s
}
