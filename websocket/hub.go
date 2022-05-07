package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//Will centrilize all the clients

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//Define the structure for the hub
type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

//Constructor for the new hub
func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

//Define the new route
func (hub *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	client := NewClient(hub, socket)
	hub.register <- client
	//Write a go routine
	go client.Write()
}

func (hub *Hub) Run() {
	for {
		//Multiplex
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

func (hub *Hub) onConnect(client *Client) {
	log.Println("Client Connected", client.socket.RemoteAddr())
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	client.id = client.socket.RemoteAddr().String()
	hub.clients = append(hub.clients, client)
}

func (hub *Hub) onDisconnect(client *Client) {
	log.Println("Client Disconnected", client.socket.RemoteAddr())
	client.socket.Close()
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	i := -1
	//iterate the clients to know which client has been disconnected
	for j, c := range hub.clients {
		if c.id == client.id {
			i = j
		}
	}
	copy(hub.clients[i:], hub.clients[i+1:])

	hub.clients[len(hub.clients)-1] = nil
	//Erase from the array
	hub.clients = hub.clients[:len(hub.clients)-1]
}

//Function to
func (hub *Hub) Broadcast(message interface{}, ignore *Client) {
	//Serialize the message
	data, _ := json.Marshal(message)
	for _, client := range hub.clients {
		if client != ignore {
			client.outbound <- data
		}
	}
}
