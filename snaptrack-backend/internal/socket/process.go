package socket

import (
	"encoding/json"
	"log"
	"sort"
	"time"

	"github.com/shirou/gopsutil/process"
)

type RunningProcess struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`
	CPUPercent float64 `json:"cpu_percent"`
	MemPercent float64 `json:"mem_percent"`
	Status     string  `json:"status"`
	ReadBytes  uint64  `json:"read_bytes"`
	WriteBytes uint64  `json:"write_bytes"`
}

// MonitorAndBroadcastRunningProcesses filters top N processes and sends to WebSocket
func MonitorAndBroadcastRunningProcesses(broadcast chan<- []byte, interval time.Duration, topN int) {
	for {
		processes, err := process.Processes()
		if err != nil {
			log.Printf("Error fetching processes: %v", err)
			time.Sleep(interval)
			continue
		}

		var procList []RunningProcess

		for _, p := range processes {
			name, err := p.Name()
			if err != nil {
				continue
			}

			status, err := p.Status()
			if err != nil {
				status = "unknown"
			}

			cpu, err := p.CPUPercent()
			if err != nil {
				cpu = 0
			}

			mem, err := p.MemoryPercent()
			if err != nil {
				mem = 0
			}

			// Skip idle or zombie processes
			if cpu < 0.1 && mem < 0.1 {
				continue
			}

			ioCounters, err := p.IOCounters()
			var readBytes uint64
			var writeBytes uint64
			if err == nil && ioCounters != nil {
				readBytes = ioCounters.ReadBytes
				writeBytes = ioCounters.WriteBytes
			}

			procList = append(procList, RunningProcess{
				PID:        p.Pid,
				Name:       name,
				CPUPercent: cpu,
				MemPercent: float64(mem),
				Status:     status,
				ReadBytes:  readBytes,
				WriteBytes: writeBytes,
			})
		}

		// Sort by CPU usage descending
		sort.Slice(procList, func(i, j int) bool {
			return procList[i].CPUPercent > procList[j].CPUPercent
		})

		// Limit to top N
		if len(procList) > topN {
			procList = procList[:topN]
		}

		payload := map[string]interface{}{
			"type":      "process_list",
			"processes": procList,
		}

		data, err := json.Marshal(payload)
		if err != nil {
			log.Println("Error marshalling process data:", err)
		} else {
			broadcast <- data
		}

		time.Sleep(interval)
	}
}
