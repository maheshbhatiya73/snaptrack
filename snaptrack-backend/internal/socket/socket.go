package socket

import (
	"log"
	"net/http"
	"sync"
	"time"
	"encoding/json"


	"github.com/gorilla/websocket"
)

type ClientManager struct {
	clients   map[*websocket.Conn]bool
	lock      sync.Mutex
	broadcast chan []byte
}

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

func StartWebSocketServer() http.Handler {
    manager := NewClientManager()
    go manager.Start()

    // Note: Removed the goroutine for MonitorAndBroadcastSystemServices
    // It will be called per client connection

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Printf("WebSocket upgrade error: %v", err)
            http.Error(w, "Could not upgrade to WebSocket", http.StatusBadRequest)
            return
        }
        manager.AddClient(conn)
        log.Printf("WebSocket connected: %s", conn.RemoteAddr().String())

        // Fetch and send services data once upon connection
	go MonitorAndBroadcastSystemStats(manager.broadcast, 2*time.Second)

        go MonitorAndBroadcastSystemServices(manager.broadcast)

        for {
            messageType, message, err := conn.ReadMessage()
            if err != nil {
                log.Printf("WebSocket disconnected: %s, reason: %v", conn.RemoteAddr().String(), err)
                manager.RemoveClient(conn)
                break
            }
            // Handle incoming messages (e.g., actions)
            if messageType == websocket.TextMessage {
                var action struct {
                    Type    string `json:"type"`
                    Service string `json:"service"`
                }
                if err := json.Unmarshal(message, &action); err == nil {
                    ServicesActions(conn, action.Type, action.Service)
                }
            }
        }
    })
}