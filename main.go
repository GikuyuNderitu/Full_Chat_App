package main

import (
	"log"
	"net/http"
	"sync"

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

	// Handle Serving Static Files
	assetFiles := http.FileServer(http.Dir("assets"))
	distFiles := http.FileServer(http.Dir("dist"))
	bowerFiles := http.FileServer(http.Dir("bower_components"))
	srcFiles := http.FileServer(http.Dir("src"))

	go handleMessages()

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/ws", handleConnections)
	http.Handle("/assets/", http.StripPrefix("/assets/", assetFiles))
	http.Handle("/dist/", http.StripPrefix("/dist/", distFiles))
	http.Handle("/bower_components/", http.StripPrefix("/bower_components/", bowerFiles))
	http.Handle("/src/", http.StripPrefix("/src/", srcFiles))
	log.Println("http server started on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	wg.Wait()
}

// HomeHandler serves the index.html file
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
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
