package route

import (
	"net/http"
	"snaptrackserver/internal/controller"
)

func RegisterBackupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/backups", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateBackup(w, r)
		case http.MethodGet:
			controller.GetAllBackups(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/backups/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetBackupByID(w, r)
		case http.MethodPut:
			controller.UpdateBackup(w, r)
		case http.MethodDelete:
			controller.DeleteBackup(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
