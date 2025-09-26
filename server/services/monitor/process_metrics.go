package monitor

import (
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// readProcessesLocal collects local process information.
func readProcessesLocal() ([]ProcessInfo, error) {
	cmd := exec.Command("ps", "-eo", "pid,user,%cpu,%mem,rss,vsz,stat,start_time,cmd", "--no-headers")
	out, err := cmd.Output()
	if err != nil {
		log.Debugf("Failed to read local processes: %v", err)
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var procs []ProcessInfo
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		pid, _ := strconv.Atoi(fields[0])
		cpu, _ := strconv.ParseFloat(fields[2], 64)
		mem, _ := strconv.ParseFloat(fields[3], 64)
		rss, _ := strconv.ParseUint(fields[4], 10, 64)
		vsz, _ := strconv.ParseUint(fields[5], 10, 64)
		rss *= 1024 // Convert KB to bytes
		vsz *= 1024
		fullCmd := strings.Join(fields[8:], " ")
		procs = append(procs, ProcessInfo{
			PID:       pid,
			User:      fields[1],
			CPU:       cpu,
			Mem:       mem,
			RSS:       rss,
			VMS:       vsz,
			State:     fields[6],
			StartTime: fields[7],
			Command:   fields[8],
			FullCmd:   fullCmd,
		})
	}
	sort.Slice(procs, func(i, j int) bool { return procs[i].CPU > procs[j].CPU })
	if len(procs) > 50 {
		procs = procs[:50]
	}
	return procs, nil
}