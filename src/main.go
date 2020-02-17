package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/marceaudavid/learn-go/src/db"
	"github.com/marceaudavid/learn-go/src/models"
	"github.com/marceaudavid/learn-go/src/routes"
	uuid "github.com/satori/go.uuid"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan models.Message)    // broadcast channel
var tickets = make(map[string]bool)

var upgrader = websocket.Upgrader{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Pour checker si premier msg checker si il se trouve dans la liste de client si non demamder wsTicket

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if !clients[ws] {
			if !tickets[msg.Message] {
				ws.WriteJSON("Invalid ticket, please retry later")
				ws.Close()
				break
			} else {
				delete(tickets, msg.Message)
				clients[ws] = true
			}
		}
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
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

// WebsocketTicket ...
func WebsocketTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		_, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "No valid token", http.StatusUnauthorized)
			return
		}
		ticket := uuid.NewV4().String()
		tickets[ticket] = true
		fmt.Fprintf(w, ticket)
	}
}

// LoadEnv ...
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	LoadEnv()
	db.Reset()

	go handleMessages()

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/wsTicket", WebsocketTicket)
	http.HandleFunc("/login", routes.Login)
	http.HandleFunc("/register", routes.Register)
	http.HandleFunc("/save", routes.Save)
	http.HandleFunc("/load", routes.Load)
	http.HandleFunc("/logout", routes.Logout)

	http.ListenAndServe(":1337", nil)
}
