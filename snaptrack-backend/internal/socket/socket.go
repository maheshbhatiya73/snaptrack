package socket

import (
	"log"
	"net/http"
	"sync"
	"time"
	"encoding/json"
	"strings"
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

    // Start monitoring goroutines once, when the server starts
    go MonitorAndBroadcastSystemStats(manager.broadcast, 2*time.Second)
    go MonitorAndBroadcastFirewallRules(manager.broadcast)
    go MonitorAndBroadcastRunningPorts(manager.broadcast)
    go MonitorAndBroadcastSystemServices(manager.broadcast)

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Printf("WebSocket upgrade error: %v", err)
            http.Error(w, "Could not upgrade to WebSocket", http.StatusBadRequest)
            return
        }
        manager.AddClient(conn)
        log.Printf("WebSocket connected: %s", conn.RemoteAddr().String())

        for {
            messageType, message, err := conn.ReadMessage()
            if err != nil {
                log.Printf("WebSocket disconnected: %s, reason: %v", conn.RemoteAddr().String(), err)
                manager.RemoveClient(conn)
                break
            }
            if messageType == websocket.TextMessage {
                var action struct {
                    Type string          `json:"type"`
                    Data json.RawMessage `json:"data"`
                }
                if err := json.Unmarshal(message, &action); err == nil {
                    HandleWebSocketAction(conn, action.Type, action.Data)
                }
            }
        }
    })
}

func HandleWebSocketAction(conn *websocket.Conn, actionType string, data json.RawMessage) {
	switch {
	case strings.HasPrefix(actionType, "services_"):
		// Extract the actual action (e.g., "start" from "services_start")
		action := strings.TrimPrefix(actionType, "services_")
		var serviceName string
		if err := json.Unmarshal(data, &serviceName); err != nil {
			sendError(conn, "Invalid service data: "+err.Error())
			return
		}
		ServicesActions(conn, action, serviceName)

	case strings.HasPrefix(actionType, "firewall_") || actionType == "stop_port" || actionType == "add_port":
		// Handle firewall and port-related actions
		FirewallActions(conn, actionType, data)

	default:
		sendError(conn, "Unknown action: "+actionType)
	}
}

func sendError(conn *websocket.Conn, message string) {
	response := struct {
		Type    string `json:"type"`
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Type:    "error",
		Success: false,
		Message: message,
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling error response: %v", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Error sending error message: %v", err)
	}
}