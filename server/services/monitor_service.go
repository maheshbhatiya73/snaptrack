package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"snaptrack/db"

	"golang.org/x/crypto/ssh"
)

type ProcessInfo struct {
	PID     int     `json:"pid"`
	User    string  `json:"user"`
	CPU     float64 `json:"cpu_percent"`
	Mem     float64 `json:"mem_percent"`
	Command string  `json:"command"`
}

type ServerMetrics struct {
	Timestamp      time.Time     `json:"timestamp"`
	ServerID       uint          `json:"server_id"`
	Host           string        `json:"host"`
	UptimeSeconds  float64       `json:"uptime_seconds"`
	Load1          float64       `json:"load1"`
	Load5          float64       `json:"load5"`
	Load15         float64       `json:"load15"`
	CPUPercent     float64       `json:"cpu_percent"`
	MemTotalBytes  uint64        `json:"mem_total_bytes"`
	MemUsedBytes   uint64        `json:"mem_used_bytes"`
	DiskTotalBytes uint64        `json:"disk_total_bytes"`
	DiskUsedBytes  uint64        `json:"disk_used_bytes"`
	Processes      []ProcessInfo `json:"processes,omitempty"`
	Error          *string       `json:"error"`
}

// ---------------- Local Collectors ----------------

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

func readLoadAverages() (float64, float64, float64, error) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return 0, 0, 0, err
	}
	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return 0, 0, 0, fmt.Errorf("unexpected loadavg format")
	}
	l1, _ := strconv.ParseFloat(fields[0], 64)
	l5, _ := strconv.ParseFloat(fields[1], 64)
	l15, _ := strconv.ParseFloat(fields[2], 64)
	return l1, l5, l15, nil
}

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
				v, _ := strconv.ParseUint(fields[1], 10, 64)
				memTotal = v * 1024
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				v, _ := strconv.ParseUint(fields[1], 10, 64)
				memAvailable = v * 1024
			}
		}
	}
	if memTotal == 0 {
		return 0, 0, fmt.Errorf("memtotal 0")
	}
	used = memTotal - memAvailable
	return memTotal, used, nil
}

func readCPUPercentSample() (idle, total uint64, err error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	line := strings.SplitN(string(data), "\n", 2)[0]
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return 0, 0, fmt.Errorf("bad /proc/stat")
	}
	var nums []uint64
	for _, f := range fields[1:] {
		v, _ := strconv.ParseUint(f, 10, 64)
		nums = append(nums, v)
	}
	idle = nums[3]
	for _, v := range nums {
		total += v
	}
	return idle, total, nil
}

func readCPUPercent() (float64, error) {
	idle1, total1, err := readCPUPercentSample()
	if err != nil {
		return 0, err
	}
	time.Sleep(500 * time.Millisecond)
	idle2, total2, err := readCPUPercentSample()
	if err != nil {
		return 0, err
	}
	idleDelta := float64(idle2 - idle1)
	totalDelta := float64(total2 - total1)
	if totalDelta <= 0 {
		return 0, nil
	}
	usage := (1.0 - idleDelta/totalDelta) * 100.0
	if usage < 0 {
		usage = 0
	}
	if usage > 100 {
		usage = 100
	}
	return usage, nil
}

func readDiskRoot() (total, used uint64, err error) {
	cmd := exec.Command("df", "-PB1", "/")
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) < 2 {
		return 0, 0, fmt.Errorf("unexpected df output")
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 5 {
		return 0, 0, fmt.Errorf("unexpected df fields")
	}
	total64, _ := strconv.ParseUint(fields[1], 10, 64)
	used64, _ := strconv.ParseUint(fields[2], 10, 64)
	return total64, used64, nil
}

func readProcessesLocal() ([]ProcessInfo, error) {
	cmd := exec.Command("ps", "-eo", "pid,user,%cpu,%mem,comm", "--no-headers")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var procs []ProcessInfo
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		pid, _ := strconv.Atoi(fields[0])
		cpu, _ := strconv.ParseFloat(fields[2], 64)
		mem, _ := strconv.ParseFloat(fields[3], 64)
		procs = append(procs, ProcessInfo{
			PID: pid, User: fields[1], CPU: cpu, Mem: mem, Command: fields[4],
		})
	}
	sort.Slice(procs, func(i, j int) bool { return procs[i].CPU > procs[j].CPU })
	if len(procs) > 50 {
		procs = procs[:50]
	}
	return procs, nil
}

func getLocalMetrics(server db.Server) ServerMetrics {
	m := ServerMetrics{Timestamp: time.Now(), ServerID: server.ID, Host: server.Host, Processes: []ProcessInfo{}}
	if up, err := readFirstFloat("/proc/uptime"); err == nil {
		m.UptimeSeconds = up
	}
	if l1, l5, l15, err := readLoadAverages(); err == nil {
		m.Load1, m.Load5, m.Load15 = l1, l5, l15
	}
	if cpu, err := readCPUPercent(); err == nil {
		m.CPUPercent = cpu
	}
	if mt, mu, err := readMemInfo(); err == nil {
		m.MemTotalBytes, m.MemUsedBytes = mt, mu
	}
	if dt, du, err := readDiskRoot(); err == nil {
		m.DiskTotalBytes, m.DiskUsedBytes = dt, du
	}
	if procs, err := readProcessesLocal(); err == nil {
		m.Processes = procs
	}
	return m
}

// ---------------- Remote Collectors ----------------

func sshClient(server db.Server) (*ssh.Client, error) {
	if server.SSHKeyPath == nil || server.SSHUser == nil || server.SSHPort == nil {
		return nil, fmt.Errorf("missing ssh config")
	}
	key, err := os.ReadFile(*server.SSHKeyPath)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User:            *server.SSHUser,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	addr := net.JoinHostPort(server.Host, strconv.Itoa(*server.SSHPort))
	return ssh.Dial("tcp", addr, config)
}

func runSSHCommand(client *ssh.Client, cmd string) (string, error) {
	fmt.Printf("[DEBUG] Running SSH command: %s\n", cmd)
	s, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer s.Close()
	out, err := s.CombinedOutput("bash -lc \"" + cmd + "\"")
	fmt.Printf("[DEBUG] Output: %s\n", string(out))
	fmt.Printf("[DEBUG] Error: %v\n", err)
	return string(out), err
}

func readProcessesRemote(client *ssh.Client) ([]ProcessInfo, error) {
	out, err := runSSHCommand(client, "ps -eo pid,user,%cpu,%mem,comm --no-headers")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	var procs []ProcessInfo
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		pid, _ := strconv.Atoi(fields[0])
		cpu, _ := strconv.ParseFloat(fields[2], 64)
		mem, _ := strconv.ParseFloat(fields[3], 64)
		procs = append(procs, ProcessInfo{
			PID: pid, User: fields[1], CPU: cpu, Mem: mem, Command: fields[4],
		})
	}
	sort.Slice(procs, func(i, j int) bool { return procs[i].CPU > procs[j].CPU })
	if len(procs) > 50 {
		procs = procs[:50]
	}
	return procs, nil
}

func getRemoteMetrics(server db.Server) ServerMetrics {
	m := ServerMetrics{Timestamp: time.Now(), ServerID: server.ID, Host: server.Host, Processes: []ProcessInfo{}}
	client, err := sshClient(server)
	if err != nil {
		msg := err.Error()
		m.Error = &msg
		return m
	}
	defer client.Close()

	// --- Uptime ---
	if out, err := runSSHCommand(client, "cut -d\\  -f1 /proc/uptime"); err == nil {
		if v, err := strconv.ParseFloat(strings.TrimSpace(out), 64); err == nil {
			m.UptimeSeconds = v
		}
	}

	// --- Loadavg ---
	if out, err := runSSHCommand(client, "cut -d' ' -f1-3 /proc/loadavg"); err == nil {
		f := strings.Fields(out)
		if len(f) >= 3 {
			m.Load1, _ = strconv.ParseFloat(f[0], 64)
			m.Load5, _ = strconv.ParseFloat(f[1], 64)
			m.Load15, _ = strconv.ParseFloat(f[2], 64)
		}
	}

	// --- Memory ---
	if out, err := runSSHCommand(client, `grep -E "MemTotal|MemAvailable" /proc/meminfo | awk '{print $2}'`); err == nil {
		lines := strings.Fields(out)
		if len(lines) >= 2 {
			tot, _ := strconv.ParseUint(lines[0], 10, 64)
			avail, _ := strconv.ParseUint(lines[1], 10, 64)
			m.MemTotalBytes = tot * 1024
			m.MemUsedBytes = (tot - avail) * 1024
		}
	}

	// --- Disk ---
	if out, err := runSSHCommand(client, `df -PB1 / | tail -1 | awk '{print $2,$3}'`); err == nil {
		fields := strings.Fields(out)
		if len(fields) >= 2 {
			m.DiskTotalBytes, _ = strconv.ParseUint(fields[0], 10, 64)
			m.DiskUsedBytes, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}

	// --- CPU ---
	if out, err := runSSHCommand(client, `
read -r c a b c d e f g h i < /proc/stat
idle1=$d
total1=0
for v in $a $b $c $d $e $f $g $h $i; do total1=$((total1+v)); done
sleep 0.5
read -r c a b c d e f g h i < /proc/stat
idle2=$d
total2=0
for v in $a $b $c $d $e $f $g $h $i; do total2=$((total2+v)); done
idle=$((idle2-idle1))
total=$((total2-total1))
if [ $total -gt 0 ]; then
awk -v i=$idle -v t=$total 'BEGIN{printf "%.2f\n", (1 - i/t)*100}'; else echo 0; fi
`); err == nil {
		if v, err := strconv.ParseFloat(strings.TrimSpace(out), 64); err == nil {
			m.CPUPercent = v
		}
	}

	// --- Processes ---
	if procs, err := readProcessesRemote(client); err == nil {
		m.Processes = procs
	}

	return m
}

// ---------------- Public API ----------------

func CollectServerMetrics(server db.Server) ServerMetrics {
	if server.Type == "remote" {
		return getRemoteMetrics(server)
	}
	return getLocalMetrics(server)
}

func (m ServerMetrics) JSON() []byte {
	b, _ := json.Marshal(m)
	return b
}
