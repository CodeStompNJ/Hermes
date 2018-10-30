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

type Todo struct {
	Name      string
	Completed bool
	Num		int
}

type Todos []Todo

var clients = make(map[*websocket.Conn]bool) //connected clients
var broadcast = make(chan MessageTest)           //broadcast channel

type MessageTest struct {
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

func SetupRouting() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", HandleConnections)
	// group routing
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

		//Probably should deal with regex outside of server.go
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

func GroupHistory(w http.ResponseWriter, r *http.Request) {
	todos := Todos{
		Todo{Name: "Write presentation", Num: 1},
		Todo{Name: "Host meetup", Num: 2},
	}

	/*sample2 := pg.Messages{
		pg.Message{ID: 1, text: "Pog"},
		pg.Message{ID: 2, text: "Champ" },
	}
	sample3 := pg.Message{ID: 1, text: "Hello"}*/
	//pass in group ID from somewhere
	fmt.Printf("LUL")
	sample := pg.GetMessagesForRoom(1)
	//fmt.Printf(string(sample))
	fmt.Printf("%v", sample)

	json.NewEncoder(w).Encode(todos)
	json.NewEncoder(w).Encode(sample)
	//json.NewEncoder(w).Encode(sample3)

	fmt.Printf("DOING IT")

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
