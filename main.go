package main

import (
	"log"
	"net/http"
	"sync"

	mgo "gopkg.in/mgo.v2"

	"github.com/GikuyuNderitu/chat_application/server/controllers"
	"github.com/gorilla/websocket"
)

var wg = sync.WaitGroup{}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}

// Message is a struct sent from a user
type Message struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	wg.Add(1)

	chatDB := createConnection("chat")
	uc := controllers.NewUserController(chatDB)

	// r := mux.NewRouter()
	// r.NotFoundHandler = http.HandlerFunc(HomeHandler)
	// r.HandleFunc("/", HomeHandler)
	// r.HandleFunc("/ws", handleConnections)
	// r.HandleFunc("/login", handleLogin)

	// Handle Serving Static Files
	assetFiles := http.FileServer(http.Dir("assets"))
	distFiles := http.FileServer(http.Dir("dist"))
	bowerFiles := http.FileServer(http.Dir("bower_components"))
	nodeModules := http.FileServer(http.Dir("node_modules"))
	srcFiles := http.FileServer(http.Dir("src"))
	imageFiles := http.FileServer(http.Dir("images"))

	go handleMessages()

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/register", uc.Register)
	http.HandleFunc("/login", uc.Login)
	http.HandleFunc("/delete", uc.Delete)
	http.Handle("/assets/", http.StripPrefix("/assets/", assetFiles))
	http.Handle("/dist/", http.StripPrefix("/dist/", distFiles))
	http.Handle("/bower_components/", http.StripPrefix("/bower_components/", bowerFiles))
	http.Handle("/node_modules/", http.StripPrefix("/node_modules/", nodeModules))
	http.Handle("/src/", http.StripPrefix("/src/", srcFiles))
	http.Handle("/images/", http.StripPrefix("/images/", imageFiles))
	log.Println("http server started on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	wg.Wait()
}

func createConnection(endpoint string) *mgo.Database {
	// Connect to local mongodb
	connection, err := mgo.Dial("mongodb://localhost")

	// Check if connection error, is mongo is running
	if err != nil {
		panic(err)
	}
	return connection.DB(endpoint)
}

// HomeHandler serves the index.html file
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {

}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	clients[ws] = true
	// log.Printf("New Connection: %v\n", ws)

	for {
		var msg Message

		log.Printf("Entered for")

		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Connection error: %v\n", err)
			delete(clients, ws)
			break
		}

		log.Printf("New Message: %v\n", msg)

		broadcast <- msg
	}
}

func handleMessages() {
	log.Printf("From HandleMessages:\n")
	for {
		log.Printf("Started Blocking:\n")
		msg := <-broadcast

		log.Printf("From HandleMessages: %v\n", msg)

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Message broadcast error: %v\n", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
	wg.Done()
}
