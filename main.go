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

	http.HandleFunc("/ws", server.HandleConnections)
	http.HandleFunc("/history", server.GroupHistory)
	http.HandleFunc("/user", server.ShowUser)

	//start listening for incoming chat messages
	go server.HandleMessages()

	//start the server on local host port 8000 and log any errors
	log.Println("https server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
