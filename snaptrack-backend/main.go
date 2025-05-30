package main

import (
	"fmt"
	"log"
	"net/http"
	"snaptrackserver/internal/config"
	"snaptrackserver/internal/route"
	"snaptrackserver/internal/api"
	"snaptrackserver/internal/socket"
	"snaptrackserver/internal/controller"
	"snaptrackserver/internal/services"

)

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.Header.Get("Origin"))
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE , OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			w.Header().Set("Content-Type", "application/json")
		} else {
			log.Printf("CORS rejected: invalid origin %s", origin)
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
	
	services.InitBackupService(db)
    services.PrintAllBackupJobs()
    controller.SetCollection(db.Collection("backups"))
	addr := fmt.Sprintf(":%d", port)
	webSocketServer := socket.StartWebSocketServer()

	mux := http.NewServeMux()
	mux.Handle("/ws", webSocketServer)
	mux.HandleFunc("/api/login", api.LoginHandler)
	mux.HandleFunc("/api/verify", api.VerifyHandler)
	route.RegisterBackupRoutes(mux) // ← Add this

	log.Printf("Server started at %s", addr)
	return http.ListenAndServe(addr, corsMiddleware(mux))
}


func main() {
	if err := Start(8000); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}