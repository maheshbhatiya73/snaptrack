package socket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"fmt"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemStats struct {
	DiskTotal   uint64  `json:"diskTotal"`
	DiskUsed    uint64  `json:"diskUsed"`
	CPUPercent  float64 `json:"cpuPercent"`
	MemTotal    uint64  `json:"ramTotalBytes"`
	MemUsed     uint64  `json:"ramUsedBytes"`
	Uptime      uint64  `json:"uptimeSeconds"`
	NetworkSent uint64  `json:"netOutBytes"`
	NetworkRecv uint64  `json:"netInBytes"`
	RAMPercent  float64 `json:"ramPercent"`
	DiskPercent float64 `json:"diskPercent"`
}


func GetDiskUsage() (total uint64, used uint64, err error) {
    partitions, err := disk.Partitions(false)
    if err != nil {
        return 0, 0, err
    }
    if len(partitions) == 0 {
        return 0, 0, fmt.Errorf("no partitions found")
    }

    validFsTypes := map[string]bool{
        "ext4":  true,
        "xfs":   true,
        "btrfs": true,
        "ntfs":  true,
        "apfs":  true,
        "hfs":   true,
    }

    for _, p := range partitions {
        if validFsTypes[p.Fstype] {
            usage, err := disk.Usage(p.Mountpoint)
            if err != nil {
                log.Printf("Error getting usage for %s: %v", p.Mountpoint, err)
                continue
            }
            if usage.Total > 0 {
                return usage.Total, usage.Used, nil
            }
        }
    }

    return 0, 0, fmt.Errorf("no valid disk usage info found")
}

func GetCPUUsage() (float64, error) {
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) == 0 {
		return 0, nil
	}
	return percentages[0], nil
}

func GetRAMUsage() (total uint64, used uint64, err error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, err
	}
	return vmStat.Total, vmStat.Used, nil
}

func GetUptime() (uint64, error) {
	return host.Uptime()
}

func GetNetworkUsage() (sent uint64, recv uint64, err error) {
	netIOCounters, err := net.IOCounters(false)
	if err != nil || len(netIOCounters) == 0 {
		return 0, 0, err
	}
	return netIOCounters[0].BytesSent, netIOCounters[0].BytesRecv, nil
}


func FetchStats() (*SystemStats, error) {
	diskTotal, diskUsed, err := GetDiskUsage()
	if err != nil {
		log.Printf("Error getting disk usage: %v", err)
	}
	cpuPercent, err := GetCPUUsage()
	if err != nil {
		log.Printf("Error getting CPU usage: %v", err)
	}
	memTotal, memUsed, err := GetRAMUsage()
	if err != nil {
		log.Printf("Error getting RAM usage: %v", err)
	}
	uptime, err := GetUptime()
	if err != nil {
		log.Printf("Error getting uptime: %v", err)
	}
	netSent, netRecv, err := GetNetworkUsage()
	if err != nil {
		log.Printf("Error getting network usage: %v", err)
	}
	
	var ramPercent float64
	var diskPercent float64

	if memTotal > 0 {
		ramPercent = float64(memUsed) / float64(memTotal) * 100
	}

	if diskTotal > 0 {
		diskPercent = float64(diskUsed) / float64(diskTotal) * 100
	}

	return &SystemStats{
		DiskTotal:   diskTotal,
		DiskUsed:    diskUsed,
		CPUPercent:  cpuPercent,
		MemTotal:    memTotal,
		MemUsed:     memUsed,
		Uptime:      uptime,
		NetworkSent: netSent,
		NetworkRecv: netRecv,
		RAMPercent:  ramPercent,
		DiskPercent: diskPercent,
	}, nil
}


func MonitorAndBroadcast(broadcast chan []byte, interval time.Duration) {
	for {
		stats, err := FetchStats()
		if err != nil {
			log.Printf("Error fetching stats: %v", err)
			time.Sleep(interval)
			continue
		}

		data, err := json.Marshal(stats)
		if err != nil {
			log.Printf("Error marshalling stats: %v", err)
			time.Sleep(interval)
			continue
		}

		broadcast <- data
		time.Sleep(interval)
	}
}
