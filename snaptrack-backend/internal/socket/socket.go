package socket

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ClientManager manages all active websocket clients and broadcasts messages
type ClientManager struct {
	clients   map[*websocket.Conn]bool
	lock      sync.Mutex
	broadcast chan []byte
}

// NewClientManager initializes a ClientManager
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

func (cm *ClientManager) AddClient(conn *websocket.Conn) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.clients[conn] = true
}

func (cm *ClientManager) RemoveClient(conn *websocket.Conn) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	delete(cm.clients, conn)
	conn.Close()
}

// Start listens for broadcast messages and sends them to all clients
func (cm *ClientManager) Start() {
	for {
		msg := <-cm.broadcast
		cm.lock.Lock()
		for client := range cm.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Error writing to client %v: %v", client.RemoteAddr(), err)
				client.Close()
				delete(cm.clients, client)
			}
		}
		cm.lock.Unlock()
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Printf("Origin: %s", origin)
		return origin == "http://localhost:3000"
	},
}

// StartWebSocketServer returns an http.Handler that upgrades HTTP to WebSocket and sends stats
func StartWebSocketServer() http.Handler {
	manager := NewClientManager()
	go manager.Start()

	// Start system stats monitor and broadcast every 2 seconds
	go MonitorAndBroadcast(manager.broadcast, 2*time.Second)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			http.Error(w, "Could not upgrade to WebSocket", http.StatusBadRequest)
			return
		}
		manager.AddClient(conn)
		log.Printf("WebSocket connected: %s", conn.RemoteAddr().String())

		// Keep connection alive; disconnect on error
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Printf("WebSocket disconnected: %s, reason: %v", conn.RemoteAddr().String(), err)
				manager.RemoveClient(conn)
				break
			}
		}
	})
}
