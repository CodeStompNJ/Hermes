package processing

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/go-playground/validator.v9"
	"github.com/dgrijalva/jwt-go"

	pg "../postgres"

	"github.com/gorilla/websocket"
	//"github.com/gorilla/mux"
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

var clients = make(map[*websocket.Conn]bool) //connected clients
var broadcast = make(chan MessageTest)       //broadcast channel

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

	//JWT endpoints
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)

	//registration
	http.HandleFunc("/register", Register)
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
	log.Println("creating new message")
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

// Create the Signin handler
func Signin(w http.ResponseWriter, r *http.Request) {
	AddCors(&w)
	fmt.Printf("got here lul")
	log.Println("omegalul")
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
		fmt.Println("Validation failed:", e)
	} else {
		// right now not sending firstname and lastname
		// need to decide if it's something we want to keep or cut from the registerUser struct
		pg.CreateUser(t.Username,"firstname","lastname",t.Email,t.Password)
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
