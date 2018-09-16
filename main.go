package main

import (
	"log"
	"net/http"

	pg "./postgres"

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
	Group    string `json:"group"`
}

func main() {

	demo()

	pg.OpenDBConnection()
	pg.SetupDB()
	defer pg.CloseDBConnection()

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
