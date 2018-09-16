package processing

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	//"github.com/gorilla/mux"
	//"github.com/gorilla/websocket"
)

type Todo struct {
	Name      string
	Completed bool
}

type Todos []Todo

var clients = make(map[*websocket.Conn]bool) //connected clients
var broadcast = make(chan Message)           //broadcast channel

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
	Group    string `json:"group"`
}

//configure the upgrader
var upgrader = websocket.Upgrader{}

func Demo() {
	fmt.Println("READING FRMO ANOTHER FILE")
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	//Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	//Close the connection when the function returns
	defer ws.Close()

	//Register clients
	clients[ws] = true

	for {
		var msg Message

		//Read in a new message as JSON and map it to the Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		//Grab the next message from the broadcast channel

		msg := <-broadcast

		regExMesg := msg.Message

		cmds := []string{"!!age;", "!!name;", "!!hello;"}

		for _, cmds := range cmds {
			regExMesg = replaceCommands(regExMesg, cmds)
		}

		msg.Message = regExMesg

		//Send it to every client that is connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func SampleHistory(w http.ResponseWriter, r *http.Request) {
	todos := Todos{
		Todo{Name: "Write presentation"},
		Todo{Name: "Host meetup"},
	}

	json.NewEncoder(w).Encode(todos)
}