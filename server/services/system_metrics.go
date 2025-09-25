package services

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"snaptrack/db"
	"strconv"
	"strings"
	"time"
)

// readFirstFloat reads the first float value from a file.
func readFirstFloat(path string) (float64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return 0, fmt.Errorf("no data in %s", path)
	}
	return strconv.ParseFloat(fields[0], 64)
}

// detectContainer checks if the system is running in a container.
func detectContainer() bool {
	_, err := os.Stat("/.dockerenv")
	if err == nil {
		return true
	}
	data, err := os.ReadFile("/proc/1/cgroup")
	if err == nil && (strings.Contains(string(data), "docker") || strings.Contains(string(data), "lxc")) {
		return true
	}
	return false
}

// readCPUCores returns the number of CPU cores.
func readCPUCores() int {
	cmd := exec.Command("nproc", "--all")
	out, err := cmd.Output()
	if err != nil {
		log.Debugf("Failed to read CPU cores: %v", err)
		return 0
	}
	cores, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	return cores
}

// readLoggedInUsers returns the number of logged-in users.
func readLoggedInUsers() int {
	cmd := exec.Command("who", "-q")
	out, err := cmd.Output()
	if err != nil {
		log.Debugf("Failed to read logged-in users: %v", err)
		return 0
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 1 {
		count, _ := strconv.Atoi(strings.TrimSpace(lines[1]))
		return count
	}
	return 0
}

// readSystemTemp reads the system temperature in Celsius.
func readSystemTemp() float64 {
	data, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		log.Debugf("Failed to read system temperature: %v", err)
		return 0
	}
	temp, _ := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	return temp / 1000 // Convert millidegrees to Celsius
}

// readSwapInfo returns total and used swap memory in bytes.
func readSwapInfo() (total, used uint64, err error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	var swapTotal, swapFree uint64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "SwapTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				v, _ := strconv.ParseUint(fields[1], 10, 64)
				swapTotal = v * 1024
			}
		} else if strings.HasPrefix(line, "SwapFree:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				v, _ := strconv.ParseUint(fields[1], 10, 64)
				swapFree = v * 1024
			}
		}
	}
	used = swapTotal - swapFree
	return swapTotal, used, nil
}

// readDiskIO returns disk read and write bytes.
func readDiskIO() (read, write uint64, err error) {
	data, err := os.ReadFile("/proc/diskstats")
	if err != nil {
		return 0, 0, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 10 && (fields[2] == "xvda" || fields[2] == "sda") {
			read, _ = strconv.ParseUint(fields[5], 10, 64)  // Sectors read
			write, _ = strconv.ParseUint(fields[9], 10, 64) // Sectors written
			read *= 512                                     // Convert sectors to bytes
			write *= 512
			return read, write, nil
		}
	}
	return 0, 0, fmt.Errorf("no disk stats found")
}

// readNetIO returns network sent and received bytes.
func readNetIO() (sent, recv uint64, err error) {
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return 0, 0, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines[2:] { // Skip headers
		fields := strings.Fields(strings.TrimLeft(line, " "))
		if len(fields) > 16 || strings.HasPrefix(fields[0], "lo:") {
			continue
		}
		r, _ := strconv.ParseUint(fields[1], 10, 64) // Recv bytes
		s, _ := strconv.ParseUint(fields[9], 10, 64) // Sent bytes
		recv += r
		sent += s
	}
	return sent, recv, nil
}

// readLoadAverages reads system load averages (1, 5, 15 minutes).
func readLoadAverages() (load1, load5, load15 float64, err error) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return 0, 0, 0, err
	}
	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return 0, 0, 0, fmt.Errorf("insufficient data in /proc/loadavg")
	}
	load1, _ = strconv.ParseFloat(fields[0], 64)
	load5, _ = strconv.ParseFloat(fields[1], 64)
	load15, _ = strconv.ParseFloat(fields[2], 64)
	return load1, load5, load15, nil
}

// readCPUPercent reads the CPU usage percentage.
func readCPUPercent() (float64, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0, err
	}
	defer f.Close()
	var user, nice, system, idle, iowait, irq, softirq, steal, guest, guestNice uint64
	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			user, _ = strconv.ParseUint(fields[1], 10, 64)
			nice, _ = strconv.ParseUint(fields[2], 10, 64)
			system, _ = strconv.ParseUint(fields[3], 10, 64)
			idle, _ = strconv.ParseUint(fields[4], 10, 64)
			iowait, _ = strconv.ParseUint(fields[5], 10, 64)
			irq, _ = strconv.ParseUint(fields[6], 10, 64)
			softirq, _ = strconv.ParseUint(fields[7], 10, 64)
			steal, _ = strconv.ParseUint(fields[8], 10, 64)
			guest, _ = strconv.ParseUint(fields[9], 10, 64)
			guestNice, _ = strconv.ParseUint(fields[10], 10, 64)
		}
	}
	totalPrev := user + nice + system + idle + iowait + irq + softirq + steal + guest + guestNice
	idlePrev := idle
	time.Sleep(500 * time.Millisecond)
	f, err = os.Open("/proc/stat")
	if err != nil {
		return 0, err
	}
	defer f.Close()
	scanner = bufio.NewScanner(f)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			user, _ = strconv.ParseUint(fields[1], 10, 64)
			nice, _ = strconv.ParseUint(fields[2], 10, 64)
			system, _ = strconv.ParseUint(fields[3], 10, 64)
			idle, _ = strconv.ParseUint(fields[4], 10, 64)
			iowait, _ = strconv.ParseUint(fields[5], 10, 64)
			irq, _ = strconv.ParseUint(fields[6], 10, 64)
			softirq, _ = strconv.ParseUint(fields[7], 10, 64)
			steal, _ = strconv.ParseUint(fields[8], 10, 64)
			guest, _ = strconv.ParseUint(fields[9], 10, 64)
			guestNice, _ = strconv.ParseUint(fields[10], 10, 64)
		}
	}
	total := user + nice + system + idle + iowait + irq + softirq + steal + guest + guestNice
	totalDelta := total - totalPrev
	idleDelta := idle - idlePrev
	if totalDelta == 0 {
		return 0, nil
	}
	return (1.0 - float64(idleDelta)/float64(totalDelta)) * 100, nil
}

// readMemInfo reads total and used memory in bytes.
func readMemInfo() (total, used uint64, err error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	var memTotal, memAvailable uint64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				memTotal, _ = strconv.ParseUint(fields[1], 10, 64)
				memTotal *= 1024
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				memAvailable, _ = strconv.ParseUint(fields[1], 10, 64)
				memAvailable *= 1024
			}
		}
	}
	used = memTotal - memAvailable
	return memTotal, used, nil
}

// readDiskRoot reads total and used disk space for the root filesystem in bytes.
func readDiskRoot() (total, used uint64, err error) {
	cmd := exec.Command("df", "-B1", "/")
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return 0, 0, fmt.Errorf("no disk data")
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return 0, 0, fmt.Errorf("insufficient disk data")
	}
	total, _ = strconv.ParseUint(fields[1], 10, 64)
	used, _ = strconv.ParseUint(fields[2], 10, 64)
	return total, used, nil
}

// getLocalMetrics collects metrics for the local server.
func getLocalMetrics(server db.Server) ServerMetrics {
	m := ServerMetrics{
		Timestamp: time.Now(),
		ServerID:  server.ID,
		Host:      server.Host,
		Processes: []ProcessInfo{},
	}
	m.IsContainer = detectContainer()
	m.CPUCores = readCPUCores()
	m.LoggedInUsers = readLoggedInUsers()
	m.SystemTempCelsius = readSystemTemp()

	if up, err := readFirstFloat("/proc/uptime"); err == nil {
		m.UptimeSeconds = up
	} else {
		log.Debugf("Failed to read uptime: %v", err)
	}
	if l1, l5, l15, err := readLoadAverages(); err == nil {
		m.Load1, m.Load5, m.Load15 = l1, l5, l15
	} else {
		log.Debugf("Failed to read load averages: %v", err)
	}
	if cpu, err := readCPUPercent(); err == nil {
		m.CPUPercent = cpu
	} else {
		log.Debugf("Failed to read CPU percent: %v", err)
	}
	if mt, mu, err := readMemInfo(); err == nil {
		m.MemTotalBytes, m.MemUsedBytes = mt, mu
	} else {
		log.Debugf("Failed to read memory info: %v", err)
	}
	if st, su, err := readSwapInfo(); err == nil {
		m.SwapTotalBytes, m.SwapUsedBytes = st, su
	} else {
		log.Debugf("Failed to read swap info: %v", err)
	}
	if dt, du, err := readDiskRoot(); err == nil {
		m.DiskTotalBytes, m.DiskUsedBytes = dt, du
	} else {
		log.Debugf("Failed to read disk usage: %v", err)
	}
	if dr, dw, err := readDiskIO(); err == nil {
		m.DiskReadBytes, m.DiskWriteBytes = dr, dw
	} else {
		log.Debugf("Failed to read disk IO: %v", err)
	}
	if ns, nr, err := readNetIO(); err == nil {
		m.NetSentBytes, m.NetRecvBytes = ns, nr
	} else {
		log.Debugf("Failed to read net IO: %v", err)
	}
	if procs, err := readProcessesLocal(); err == nil {
		m.Processes = procs
	} else {
		log.Debugf("Failed to read processes: %v", err)
	}
	return m
}