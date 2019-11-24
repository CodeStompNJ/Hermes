package processing

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"

	pg "../postgres"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	//"github.com/gorilla/websocket"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// ****************************************************************************
// 								WebSocket Stuff
// ****************************************************************************
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

type StructContract struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
	Type      string `json:"type,omitempty"`
}

// global manager for websocket management
var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

type MessageTest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
	Group    string `json:"group"`
}

type registerTest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}

type resultMessage struct {
	Result string `json:"result"`
}

//configure the upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // @todo - just block webserver, and not allow everything
	},
}

// SetupRouting - set up HandleFuncs for routing services
func SetupRouting() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// http.HandleFunc("/ws", HandleConnections)

	//sends user info to front end
	http.HandleFunc("/user", ShowUser)

	// group routing, shows history
	http.HandleFunc("/history", GroupHistory)

	http.HandleFunc("/message", CreateNewMessage)

	// message routing
	http.HandleFunc("/messages", MessageHandler)

	//JWT endpoints
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)

	//registration
	http.HandleFunc("/register", Register)
}

// SetupWebSocket used to start manager
func SetupWebSocket() {
	manager.start()
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			jsonMessage, _ := json.Marshal(&StructContract{Content: "/A new socket has connected.", Type: "new_client"})
			manager.send(jsonMessage, conn)
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				jsonMessage, _ := json.Marshal(&StructContract{Content: "/A socket has disconnected.", Type: "client_leaving"})
				manager.send(jsonMessage, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				// send to all clients; includes sending to posting client.
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			break
		}
		stringMessage := getStringFromBytes(message)
		pg.CreateMessage(stringMessage, 1, 1)
		jsonMessage, _ := json.Marshal(&StructContract{Sender: c.id, Content: stringMessage, Type: "message"})
		manager.broadcast <- jsonMessage
	}
}

func getStringFromBytes(bytes []byte) string {
	if len(bytes) > 0 && bytes[0] == '"' {
		bytes = bytes[1:]
	}
	if len(bytes) > 0 && bytes[len(bytes)-1] == '"' {
		bytes = bytes[:len(bytes)-1]
	}
	stringMessage := string(bytes[:])
	return stringMessage
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func SocketMessage(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	client := &Client{id: u.String(), socket: conn, send: make(chan []byte)}

	manager.register <- client

	go client.read()
	go client.write()
}

//this is an artifact function and not currently being used as an endpoint
//however, we're keeping it for now to be used later
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

	vars := mux.Vars(r)

	groupID, ok := strconv.Atoi(vars["groupID"])

	// TODO: add error check here.
	println(ok)
	sample := pg.GetMessagesForRoom(groupID)

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
	log.Println("creating new message")
	AddCors(&w)
	decoder := json.NewDecoder(r.Body)
	var t MessageTest
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	groupID, ok := strconv.Atoi(t.Group)
	usernameID, ok := strconv.Atoi(t.Username)

	if ok != nil {
		log.Println("conversion error")
	}

	message := pg.CreateMessage(t.Message, groupID, usernameID)

	json.NewEncoder(w).Encode(message)

	//returns json encoding of the data
	pagesJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}

	fmt.Printf("%s", message)
	fmt.Printf("%s", pagesJSON)
}

// Create the Signin handler
func Signin(w http.ResponseWriter, r *http.Request) {
	AddCors(&w)
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// @TODO - need to validate creds before hitting db
	valid := pg.ValidUser(creds.Username, creds.Password)

	// Get the expected password from our in memory map
	// expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(1 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	AddCors(&w)
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Finally, return the welcome message to the user, along with their
	// username given in the token
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	AddCors(&w)

	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// (END) The code up-till this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	AddCors(&w)
	fmt.Printf("IN register")
	decoder := json.NewDecoder(r.Body)
	var t registerTest
	err := decoder.Decode(&t)
	if err != nil {
		//Panic here if there's an issue decoding the data
		panic(err)
	}
	// create a validator that we'll be using to validate our user values with
	v := validator.New()
	/* With the values we've gotten from the front end we place them in a struct
	that we'll compare out validators to. If it fails we push the error to
	the console else we create the user in the DB. @TODO have a way to return
	the error to the front end so we can express what wen wrong to the user.
	Also will need to validate that values dont already exist in DB, */
	if e := v.Struct(t); e != nil {
		for _, e := range e.(validator.ValidationErrors) {
			fmt.Print("Validation Error ")
			fmt.Println(e)
		}
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// right now not sending firstname and lastname
		// need to decide if it's something we want to keep or cut from the registerUser struct
		resultStatus := pg.CreateUser(t.Username, "firstname", "lastname", t.Email, t.Password)
		fmt.Println(resultStatus)
		//success will have boolean value if there was success or not
		var structInst resultMessage
		structInst.Result = resultStatus
		if resultStatus != "success" {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(structInst)

		//returns json encoding of the data

		//w.WriteHeader(200)
		//json.NewEncoder(w).Encode(sample)
	}

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
