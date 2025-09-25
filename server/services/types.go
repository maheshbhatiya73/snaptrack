package services

import (
	"os"
	"time"
	"github.com/sirupsen/logrus"
)

// ProcessInfo holds detailed process information.
type ProcessInfo struct {
	PID        int     `json:"pid"`
	User       string  `json:"user"`
	CPU        float64 `json:"cpu_percent"`
	Mem        float64 `json:"mem_percent"`
	RSS        uint64  `json:"rss_bytes"`    // Resident Set Size
	VMS        uint64  `json:"vms_bytes"`    // Virtual Memory Size
	State      string  `json:"state"`        // Process state (R, S, Z, etc.)
	StartTime  string  `json:"start_time"`   // Process start time
	Command    string  `json:"command"`
	FullCmd    string  `json:"full_cmdline"` // Full command line
}

// ServerMetrics holds comprehensive server metrics.
type ServerMetrics struct {
	Timestamp         time.Time     `json:"timestamp"`
	ServerID          uint          `json:"server_id"`
	Host              string        `json:"host"`
	UptimeSeconds     float64       `json:"uptime_seconds"`
	Load1             float64       `json:"load1"`
	Load5             float64       `json:"load5"`
	Load15            float64       `json:"load15"`
	CPUPercent        float64       `json:"cpu_percent"`
	CPUCores          int           `json:"cpu_cores"`
	MemTotalBytes     uint64        `json:"mem_total_bytes"`
	MemUsedBytes      uint64        `json:"mem_used_bytes"`
	SwapTotalBytes    uint64        `json:"swap_total_bytes"`
	SwapUsedBytes     uint64        `json:"swap_used_bytes"`
	DiskTotalBytes    uint64        `json:"disk_total_bytes"`
	DiskUsedBytes     uint64        `json:"disk_used_bytes"`
	DiskReadBytes     uint64        `json:"disk_read_bytes"`
	DiskWriteBytes    uint64        `json:"disk_write_bytes"`
	NetSentBytes      uint64        `json:"net_sent_bytes"`
	NetRecvBytes      uint64        `json:"net_recv_bytes"`
	LoggedInUsers     int           `json:"logged_in_users"`
	SystemTempCelsius float64       `json:"system_temp_celsius"`
	IsContainer       bool          `json:"is_container"`
	Processes         []ProcessInfo `json:"processes,omitempty"`
	Error             *string       `json:"error"`
}

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
}