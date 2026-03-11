package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	clients ClientList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

func (manager *Manager) serveWS(writer http.ResponseWriter, request *http.Request) {
	log.Println("new connection")

	// Upgrade regular http connection to websocket
	conn, err := websocketUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	var client *Client = NewClient(conn, manager)
	manager.addClient(client)
}

func (manager *Manager) addClient(client *Client) {
	manager.Lock()
	defer manager.Unlock()

	manager.clients[client] = true
}

func (manager *Manager) removeClient(client *Client) {
	manager.Lock()
	defer manager.Unlock()

	_, ok := manager.clients[client]
	if ok {
		client.connection.Close()
		delete(manager.clients, client)
	}
}

func (client *Client) readMessages() {
	defer func() {
		client.manager.removeClient(client)
	}()

	for {
		messageType, payload, err := client.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		log.Println(messageType)
		log.Println(string(payload))
	}
}
