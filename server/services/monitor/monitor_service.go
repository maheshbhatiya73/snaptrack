package monitor

import (
	"encoding/json"
	"snaptrack/db"
)

// CollectServerMetrics collects metrics for a given server based on its type (local or remote).
func CollectServerMetrics(server db.Server) ServerMetrics {
	if server.Type == "remote" {
		return getRemoteMetrics(server)
	}
	return getLocalMetrics(server)
}

// JSON serializes ServerMetrics to JSON.
func (m ServerMetrics) JSON() []byte {
	b, _ := json.Marshal(m)
	return b
}