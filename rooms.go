package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Room struct {
	id int32

	// Registered clients in current room
	clients map[*Client]bool

	// Inbound messages from the clients in room
	broadcast chan *Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// sample json message to be sent over the wire
type Message struct {
	Message  string `json:"message,omitempty"`
	Type     string `json:"type,omitempty"`
	ClientID string `json:"client_id,omitempty"`
}

// type ChatServer struct {
// 	rooms []*Room
// 	// Register requests from the clients.
// 	register chan *Room

// 	// Unregister requests from clients.
// 	unregister chan *Room
// }

func newRoom() *Room {

	// send the rand to each call to create a new rome creates a new unique ID

	rand.Seed(time.Now().UnixNano())
	room := &Room{
		id:         rand.Int31(),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}

	go room.run()
	return room
}

// this function runs an active room on the server
func (r *Room) run() {
	for {
		select {
		// registers a new client to a room
		case client := <-r.register:
			fmt.Println("client registered... room id -", client.room.id)

			r.clients[client] = true

			fmt.Println("clients", len(r.clients))
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
			fmt.Println("clients unregistered", len(r.clients))
		case message := <-r.broadcast:
			fmt.Println(message)
			for client := range r.clients {

				client.send <- message

			}
		}
	}
}
