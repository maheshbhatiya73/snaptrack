package main

import (
	"fmt"
	"log"
	"net/http"
	"snaptrackserver/internal/api"
	"snaptrackserver/internal/config"
	"snaptrackserver/internal/controller"
	"snaptrackserver/internal/route"
	"snaptrackserver/internal/services"
	"snaptrackserver/internal/socket"
)

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.Header.Get("Origin"))
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{"http://localhost:3000", "http://127.0.0.1:3000"}
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
				w.Header().Set("Content-Type", "application/json")
				break
			}
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func Start(port int) error {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	controller.SetCollection(db.Collection("backups"))
	addr := fmt.Sprintf(":%d", port)
	webSocketServer := socket.StartWebSocketServer()

	mux := http.NewServeMux()
	mux.Handle("/ws", webSocketServer)
	mux.HandleFunc("/api/login", api.LoginHandler)
	mux.HandleFunc("/api/verify", api.VerifyHandler)
	route.RegisterBackupRoutes(mux)

	log.Printf("Server started at %s", addr)
	services.InitBackupService(db)
	return http.ListenAndServe(addr, corsMiddleware(mux))
}

func main() {
	if err := Start(8000); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}