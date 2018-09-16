package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

type userInfo struct {
	Bio      string `json:"bio"`
	Age      string `json:"age"`
	Location string `json:"location"`
}

var clients = make(map[*websocket.Conn]bool) //connected clients
var broadcast = make(chan Message)           //broadcast channel

//configure the upgrader
var upgrader = websocket.Upgrader{}

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
	Group	 string `json:"group"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "seshat"
	password = "r8*W6F#8xE"
	dbname   = "hermes"
)

func main() {

	demo()

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

	defer db.Close()

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/history", sampleHistory)

	//start listening for incoming chat messages
	go handleMessages()

	//start the server on local host port 8000 and log any errors
	log.Println("https server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
