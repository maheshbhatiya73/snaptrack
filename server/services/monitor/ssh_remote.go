package monitor

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"snaptrack/db"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// sshClient establishes an SSH connection to the server.
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

	var hostKeyCallback ssh.HostKeyCallback
	knownHostsPath := os.Getenv("HOME") + "/.ssh/known_hosts"

	if _, err := os.Stat(knownHostsPath); err == nil {
		// Load original known_hosts callback
		origCallback, err := knownhosts.New(knownHostsPath)
		if err != nil {
			// log.Warnf("Failed to parse known_hosts: %v, using insecure fallback", err)
			hostKeyCallback = ssh.InsecureIgnoreHostKey()
		} else {
			// Wrap original callback to ignore mismatches safely
			hostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				if err := origCallback(hostname, remote, key); err != nil {
					// log.Warnf("Host key mismatch for %s: %v, ignoring for monitoring", hostname, err)
					return nil // ignore mismatch
				}
				return nil
			}
		}
	} else {
		// known_hosts does not exist
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	config := &ssh.ClientConfig{
		User:            *server.SSHUser,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: hostKeyCallback,
		Timeout:         10 * time.Second,
	}

	addr := net.JoinHostPort(server.Host, strconv.Itoa(*server.SSHPort))
	return ssh.Dial("tcp", addr, config)
}

// runSSHCommand executes a command over SSH and returns the output.
func runSSHCommand(client *ssh.Client, cmd string) (string, error) {
	// log.Debugf("Running SSH command: %s", cmd)
	s, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer s.Close()

	out, err := s.CombinedOutput(cmd)
	// log.Debugf("Output: %s", string(out))
	if err != nil {
		// log.Debugf("Error: %v", err)
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// readProcessesRemote collects process information from a remote server.
func readProcessesRemote(client *ssh.Client) ([]ProcessInfo, error) {
	out, err := runSSHCommand(client, "ps -eo pid,user,%cpu,%mem,rss,vsz,stat,start_time,cmd --no-headers")
	if err != nil {
		log.Debugf("Failed to read remote processes: %v", err)
		return nil, err
	}

	lines := strings.Split(out, "\n")
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
		rss *= 1024
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

// getRemoteMetrics collects metrics from a remote server via SSH.
func getRemoteMetrics(server db.Server) ServerMetrics {
	m := ServerMetrics{
		Timestamp: time.Now(),
		ServerID:  server.ID,
		Host:      server.Host,
		Processes: []ProcessInfo{},
	}

	client, err := sshClient(server)
	if err != nil {
		msg := err.Error()
		m.Error = &msg
		log.Errorf("SSH connection failed: %v", err)
		return m
	}
	defer client.Close()

	// Detect container
	if out, _ := runSSHCommand(client, "cat /.dockerenv >/dev/null 2>&1 && echo true || echo false"); out == "true" {
		m.IsContainer = true
	} else if out, _ := runSSHCommand(client, "grep -qE '(docker|lxc)' /proc/1/cgroup && echo true || echo false"); out == "true" {
		m.IsContainer = true
	}

	// CPU cores
	if out, _ := runSSHCommand(client, "nproc --all"); out != "" {
		m.CPUCores, _ = strconv.Atoi(out)
	}

	// Logged in users
	if out, _ := runSSHCommand(client, "who | wc -l"); out != "" {
		m.LoggedInUsers, _ = strconv.Atoi(out)
	}

	// System temperature
	if out, _ := runSSHCommand(client, "cat /sys/class/thermal/thermal_zone0/temp 2>/dev/null"); out != "" {
		temp, _ := strconv.ParseFloat(out, 64)
		m.SystemTempCelsius = temp / 1000
	}

	// Uptime
	if out, _ := runSSHCommand(client, "awk '{print $1}' /proc/uptime"); out != "" {
		m.UptimeSeconds, _ = strconv.ParseFloat(out, 64)
	}

	// Load average
	if out, _ := runSSHCommand(client, "awk '{print $1 \" \" $2 \" \" $3}' /proc/loadavg"); out != "" {
		fields := strings.Fields(out)
		if len(fields) >= 3 {
			m.Load1, _ = strconv.ParseFloat(fields[0], 64)
			m.Load5, _ = strconv.ParseFloat(fields[1], 64)
			m.Load15, _ = strconv.ParseFloat(fields[2], 64)
		}
	}

	// Memory usage
	if out, _ := runSSHCommand(client, "awk '/MemTotal/{t=$2} /MemAvailable/{a=$2} END{print t \" \" a}' /proc/meminfo"); out != "" {
		fields := strings.Fields(out)
		if len(fields) >= 2 {
			tot, _ := strconv.ParseUint(fields[0], 10, 64)
			avail, _ := strconv.ParseUint(fields[1], 10, 64)
			m.MemTotalBytes = tot * 1024
			m.MemUsedBytes = (tot - avail) * 1024
		}
	}

	// Swap usage
	if out, _ := runSSHCommand(client, "awk '/SwapTotal/{t=$2} /SwapFree/{f=$2} END{print t \" \" (t-f)}' /proc/meminfo"); out != "" {
		fields := strings.Fields(out)
		if len(fields) >= 2 {
			m.SwapTotalBytes, _ = strconv.ParseUint(fields[0], 10, 64)
			m.SwapUsedBytes, _ = strconv.ParseUint(fields[1], 10, 64)
			m.SwapTotalBytes *= 1024
			m.SwapUsedBytes *= 1024
		}
	}

	// Disk usage
	if out, _ := runSSHCommand(client, "df -PB1 / | tail -1 | awk '{print $2 \" \" $3}'"); out != "" {
		fields := strings.Fields(out)
		if len(fields) >= 2 {
			m.DiskTotalBytes, _ = strconv.ParseUint(fields[0], 10, 64)
			m.DiskUsedBytes, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}

	// Disk stats
	if out, _ := runSSHCommand(client, "awk '/xvda|sda/ {print $6*512 \" \" $10*512}' /proc/diskstats"); out != "" {
		fields := strings.Fields(out)
		if len(fields) >= 2 {
			m.DiskReadBytes, _ = strconv.ParseUint(fields[0], 10, 64)
			m.DiskWriteBytes, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}

	// Network stats
	if out, _ := runSSHCommand(client, "awk 'NR>2 && $1!~/lo:/ {r+=$2; s+=$10} END{print s \" \" r}' /proc/net/dev"); out != "" {
		fields := strings.Fields(out)
		if len(fields) >= 2 {
			m.NetSentBytes, _ = strconv.ParseUint(fields[0], 10, 64)
			m.NetRecvBytes, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}

	// CPU usage over 0.5s interval
	cpuCmd := `
read -r cpu user nice system idle iowait irq softirq steal guest guest_nice < /proc/stat
total=$((user + nice + system + idle + iowait + irq + softirq + steal + guest + guest_nice))
idle_prev=$idle
sleep 0.5
read -r cpu user nice system idle iowait irq softirq steal guest guest_nice < /proc/stat
total_new=$((user + nice + system + idle + iowait + irq + softirq + steal + guest + guest_nice))
idle_delta=$((idle - idle_prev))
total_delta=$((total_new - total))
if [ $total_delta -gt 0 ]; then
  usage=$(awk -v id=$idle_delta -v td=$total_delta 'BEGIN {printf "%.2f", (1 - id/td) * 100}')
else
  usage=0.00
fi
echo $usage
`
	if out, _ := runSSHCommand(client, cpuCmd); out != "" {
		m.CPUPercent, _ = strconv.ParseFloat(out, 64)
	}

	// Processes
	if procs, err := readProcessesRemote(client); err == nil {
		m.Processes = procs
	}

	return m
}
