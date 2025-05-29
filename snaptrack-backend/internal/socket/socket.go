package socket

import (
	"log"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Printf("Origin: %s", origin)
		return origin == "http://localhost:3000"
	},
}

// StartWebSocketServer starts a WebSocket server
func StartWebSocketServer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			http.Error(w, "Could not upgrade to WebSocket", http.StatusBadRequest)
			return
		}
		defer conn.Close()

		log.Printf("connected: %s", conn.RemoteAddr().String())

		// Keep connection open, no message handling
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Printf("disconnected: %s, reason: %v", conn.RemoteAddr().String(), err)
				return
			}
		}
	})
}