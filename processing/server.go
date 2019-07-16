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

	http.HandleFunc("/message", CreateNewMessage)

	// message routing
	http.HandleFunc("/messages", MessageHandler)
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	AddCors(&w)

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
	AddCors(&w)

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
	AddCors(&w)

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

func CreateNewMessage(w http.ResponseWriter, r *http.Request) {
	AddCors(&w)
	decoder := json.NewDecoder(r.Body)
	var t MessageTest
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t.Email)
	log.Println(t.Username)
	log.Println(t.Message)
	log.Println(t.Group)

	// sample := pg.CreateMessage(r.text, r.chatroomId, r.userId)
	sample := pg.CreateMessage(t.Message, 1, 1) // @TODO - need to pass in chatroom and user correctly.

	json.NewEncoder(w).Encode(sample)

	//returns json encoding of the data
	pagesJSON, err := json.Marshal(sample)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}

	fmt.Printf("%s", sample)
	fmt.Printf("%s", pagesJSON)
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	AddCors(&w)

}

func AddCors(w *http.ResponseWriter) {
	//Allow CORS here By * or specific origin
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
