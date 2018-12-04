package processing

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	pg "../postgres"

	"github.com/gorilla/websocket"
	//"github.com/gorilla/mux"
	//"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) //connected clients
var broadcast = make(chan MessageTest)       //broadcast channel

type MessageTest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
	Group    string `json:"group"`
}

//configure the upgrader
var upgrader = websocket.Upgrader{}

//set up HandleFuncs for routing services
func SetupRouting() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", HandleConnections)

	//sends user info to front end
	http.HandleFunc("/user", ShowUser)

	// group routing, shows history
	http.HandleFunc("/history", GroupHistory)

	// message routing
	http.HandleFunc("/messages", MessageHandler)
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
		var msg MessageTest

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

		//gross
		cmds := []string{"!!age;", "!!name;", "!!hello;"}

		for _, cmds := range cmds {
			regExMesg = replaceCommands(regExMesg, cmds)
		}

		//retrieve user
		//retrieve group ?

		//Probably should deal with regex outside of server.go
		msg.Message = regExMesg

		pg.CreateMessage(msg.Message, 1, 1)

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

func ShowUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("IN SHOW USER")
	//pass in user ID from somewhere, try to get from url, maybe in a form when the link is hit
	userInfo := pg.ReturnUser(1)
	fmt.Printf("IN SHOW USER2")
	//Encode writes the JSON encoding of userInfo to the stream
	json.NewEncoder(w).Encode(userInfo)
	fmt.Printf("IN SHOW USER3")
	//returns json encoding of the data
	pagesJson, err := json.Marshal(userInfo)
	fmt.Printf("IN SHOW USER4")
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}

	fmt.Printf("%s", pagesJson)
}

func GroupHistory(w http.ResponseWriter, r *http.Request) {

	//pass in group ID from somewhere
	sample := pg.GetMessagesForRoom(1)

	//fmt.Printf("%v", sample)

	json.NewEncoder(w).Encode(sample)

	pagesJson, err := json.Marshal(sample)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}
	fmt.Printf("%s", pagesJson)

	// pagesJson, err = json.Marshal(sample2)
	// if err != nil {
	//     log.Fatal("Cannot encode to JSON ", err)
	// }
	// fmt.Printf("%s", pagesJson)

	// pagesJson, err = json.Marshal(sample3)
	// if err != nil {
	//     log.Fatal("Cannot encode to JSON ", err)
	// }
	// fmt.Printf("%s", pagesJson)
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {

}
