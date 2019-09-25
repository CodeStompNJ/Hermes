package main

import (
	"log"
	"net/http"

	pg "./postgres"
	server "./processing"

	_ "github.com/lib/pq"
)

type userInfo struct {
	Bio      string `json:"bio"`
	Age      string `json:"age"`
	Location string `json:"location"`
}

func main() {

	pg.OpenDBConnection()
	pg.SetupDB()
	defer pg.CloseDBConnection()

	fs := http.FileServer(http.Dir("./public"))

	http.Handle("/", fs)
	http.HandleFunc("/ws", server.SocketMessage)
	http.HandleFunc("/history", server.GroupHistory)
	http.HandleFunc("/user", server.ShowUser)
	http.HandleFunc("/message", server.CreateNewMessage) // @TODO - make it only available to POST
	http.HandleFunc("/signin", server.Signin)
	http.HandleFunc("/welcome", server.Welcome)
	http.HandleFunc("/refresh", server.Refresh)
	http.HandleFunc("/register", server.Register)
	// server.SetupRouter()

	//start listening for incoming chat messages
	go server.SetupWebSocket()

	//start the server on local host port 8000 and log any errors
	log.Println("https server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
