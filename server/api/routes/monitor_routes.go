package routes

import (
	"encoding/json"
	"time"

	"log"
	"snaptrack/db"
	"snaptrack/services/monitor"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// GET /api/monitor/:serverID/ws
func MonitorRoutes(app *fiber.App) {
    // Batch websocket: GET /api/monitor/ws
    // Streams metrics for all servers every 2s as a JSON array
    app.Get("/api/monitor/ws", websocket.New(func(c *websocket.Conn) {
        defer c.Close()
        log.Println("[ws] /api/monitor/ws client connected")
        ticker := time.NewTicker(2 * time.Second)
        defer ticker.Stop()
        for {
            select {
            case <-ticker.C:
                var servers []db.Server
                if err := db.DB.Find(&servers).Error; err != nil {
                    // best-effort; close on error
                    return
                }
                // Collect metrics for each server
                type metricsList []monitor.ServerMetrics
                var out metricsList
                for _, s := range servers {
                    out = append(out, monitor.CollectServerMetrics(s))
                }
                if b, err := json.Marshal(out); err == nil {
                    if err := c.WriteMessage(websocket.TextMessage, b); err != nil {
                        log.Println("[ws] write error:", err)
                        return
                    }
                } else {
                    log.Println("[ws] marshal error:", err)
                    return
                }
            }
        }
    }))

	app.Get("/api/monitor/:serverID/ws", websocket.New(func(c *websocket.Conn) {
		defer c.Close()
        idParam := c.Params("serverID")
		var server db.Server
		if err := db.DB.First(&server, idParam).Error; err != nil {
			// Close silently if not found
			return
		}
        log.Println("[ws] /api/monitor/", idParam, "/ws client connected")
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				metrics := monitor.CollectServerMetrics(server)
                if err := c.WriteMessage(websocket.TextMessage, metrics.JSON()); err != nil {
                    log.Println("[ws] write error (single):", err)
					return
				}
			}
		}
	}))
}


