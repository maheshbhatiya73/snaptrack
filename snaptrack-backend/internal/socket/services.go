package socket

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type ServiceInfo struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Uptime  string `json:"uptime"`
	Memory  string `json:"memory"`
	Version string `json:"version"`
}
func MonitorAndBroadcastSystemServices(broadcast chan<- []byte) {
    for {
        services, err := GetAllSystemServices()
        if err != nil {
            logErrorWithRateLimit("Error getting services", err)
            time.Sleep(10 * time.Second)
            continue
        }

        var serviceInfos []ServiceInfo
        for _, service := range services {
            status, err := GetServiceStatus(service)
            if err != nil {
                status = "unknown"
            }

            uptime, err := GetServiceUptime(service)
            if err != nil {
                logErrorWithRateLimit(fmt.Sprintf("Error getting uptime for %s", service), err)
                uptime = "unknown"
            }

            memory, err := GetServiceMemoryUsage(service)
            if err != nil {
                logErrorWithRateLimit(fmt.Sprintf("Error getting memory for %s", service), err)
                memory = "unknown"
            }

            version, err := GetServiceVersion(service)
            if err != nil {
                logErrorWithRateLimit(fmt.Sprintf("Error getting version for %s", service), err)
                version = "n/a"
            }

            serviceInfos = append(serviceInfos, ServiceInfo{
                Name:    service,
                Status:  status,
                Uptime:  uptime,
                Memory:  memory,
                Version: version,
            })
        }

        data, err := json.Marshal(struct {
            Type     string        `json:"type"`
            Services []ServiceInfo `json:"services"`
        }{
            Type:     "services",
            Services: serviceInfos,
        })
        if err != nil {
            logErrorWithRateLimit("Error marshaling services", err)
            time.Sleep(10 * time.Second)
            continue
        }

        select {
        case broadcast <- data:
            log.Println("Sent services data to broadcast channel")
        default:
            log.Println("Broadcast channel full, skipping services update")
        }
        time.Sleep(30 * time.Second) // Poll every 30 seconds
    }
}

func GetAllSystemServices() ([]string, error) {
    services := []string{
        // Web servers
        "apache2.service",
        "httpd.service",          // CentOS/RedHat apache service name
        "nginx.service",

        // Databases
        "mysql.service",
        "mariadb.service",
        "postgresql.service",
        "mongodb.service",
        "redis.service",

        // Container runtimes
        "docker.service",
        "containerd.service",

        // SSH & remote access
        "ssh.service",
        "sshd.service",           // Some distros use sshd.service

        // System utilities
        "cron.service",
        "crond.service",          // CentOS/RedHat cron

        // Mail servers
        "postfix.service",
        "exim.service",
        "dovecot.service",

        // Caches / Message brokers
        "memcached.service",
        "rabbitmq-server.service",

        // Firewall & security
        "firewalld.service",
        "ufw.service",

        // Monitoring
        "prometheus.service",
        "node_exporter.service",

        // Logging
        "rsyslog.service",
        "syslog.service",

        // Network
        "NetworkManager.service",
        "network.service",        // Older RedHat style

        // Time sync
        "ntpd.service",
        "chronyd.service",

        // Misc
        "bluetooth.service",
        "avahi-daemon.service",
    }
    return services, nil
}

func GetServiceStatus(serviceName string) (string, error) {
	cmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := cmd.Output()
	if err != nil {
		// Instead of returning error, interpret exit code or output and send a friendly message
		// Use CombinedOutput to capture stderr (sometimes useful)
		combinedOutput, _ := exec.Command("systemctl", "is-active", serviceName).CombinedOutput()
		outStr := strings.TrimSpace(string(combinedOutput))

		// Common outputs for inactive or missing services
		switch outStr {
		case "inactive", "failed", "unknown", "activating", "deactivating":
			return outStr, nil
		case "":
			return "not available", nil
		default:
			// If unknown output or error, log but return friendly string
			return "not available", nil
		}
	}
	status := strings.TrimSpace(string(output))
	return status, nil
}
func GetServiceUptime(serviceName string) (string, error) {
	cmd := exec.Command("systemctl", "show", serviceName, "--property=ActiveEnterTimestamp")
	output, err := cmd.Output()
	if err != nil {
		return "unknown", err
	}
	timestamp := strings.TrimSpace(strings.Replace(string(output), "ActiveEnterTimestamp=", "", 1))

	// If empty or n/a, treat as "not running"
	if timestamp == "" || timestamp == "n/a" {
		return "not running", nil
	}

	layouts := []string{
		"Mon 2006-01-02 15:04:05 MST",
		"2006-01-02T15:04:05Z07:00",
	}

	var startTime time.Time
	for _, layout := range layouts {
		if t, err := time.Parse(layout, timestamp); err == nil {
			startTime = t
			break
		}
	}

	if startTime.IsZero() {
		// Instead of returning error, return a default string
		return "unknown", nil
	}

	duration := time.Since(startTime).Round(time.Second)
	return duration.String(), nil
}

func GetServiceMemoryUsage(serviceName string) (string, error) {
	cmd := exec.Command("systemctl", "show", serviceName, "--property=MemoryCurrent")
	output, err := cmd.Output()
	if err != nil {
		return "unknown", err
	}
	memory := strings.TrimSpace(strings.Replace(string(output), "MemoryCurrent=", "", 1))
	if memory == "[not set]" || memory == "" {
		return "0B", nil
	}
	bytes, err := strconv.ParseUint(memory, 10, 64)
	if err != nil {
		return "unknown", err
	}
	units := []string{"B", "KB", "MB", "GB", "TB"}
	size := float64(bytes)
	unitIndex := 0
	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}
	return fmt.Sprintf("%.2f%s", size, units[unitIndex]), nil
}
func GetServiceVersion(serviceName string) (string, error) {
	var cmdName string
	switch {
	case strings.Contains(serviceName, "apache2") || strings.Contains(serviceName, "httpd"):
		cmdName = "apache2ctl" // or "httpd" on some distros, but apache2ctl is more universal for version
	case strings.Contains(serviceName, "nginx"):
		cmdName = "nginx"
	case strings.Contains(serviceName, "mysql") || strings.Contains(serviceName, "mariadb"):
		cmdName = "mysql"
	case strings.Contains(serviceName, "postgresql"):
		cmdName = "psql"
	case strings.Contains(serviceName, "mongodb"):
		cmdName = "mongod"
	case strings.Contains(serviceName, "redis"):
		cmdName = "redis-server"
	case strings.Contains(serviceName, "ssh") || strings.Contains(serviceName, "sshd"):
		cmdName = "ssh"
	case strings.Contains(serviceName, "docker"):
		cmdName = "docker"
	case strings.Contains(serviceName, "containerd"):
		cmdName = "containerd"
	case strings.Contains(serviceName, "cron") || strings.Contains(serviceName, "crond"):
		cmdName = "crond" // or "cron" depending on system
	case strings.Contains(serviceName, "postfix"):
		cmdName = "postfix"
	case strings.Contains(serviceName, "exim"):
		cmdName = "exim"
	case strings.Contains(serviceName, "dovecot"):
		cmdName = "dovecot"
	case strings.Contains(serviceName, "memcached"):
		cmdName = "memcached"
	case strings.Contains(serviceName, "rabbitmq"):
		cmdName = "rabbitmqctl"
	case strings.Contains(serviceName, "firewalld"):
		cmdName = "firewalld"
	case strings.Contains(serviceName, "ufw"):
		cmdName = "ufw"
	case strings.Contains(serviceName, "prometheus"):
		cmdName = "prometheus"
	case strings.Contains(serviceName, "node_exporter"):
		cmdName = "node_exporter"
	case strings.Contains(serviceName, "rsyslog") || strings.Contains(serviceName, "syslog"):
		cmdName = "rsyslogd"
	case strings.Contains(serviceName, "NetworkManager"):
		cmdName = "nmcli"
	case strings.Contains(serviceName, "ntpd"):
		cmdName = "ntpd"
	case strings.Contains(serviceName, "chronyd"):
		cmdName = "chronyd"
	case strings.Contains(serviceName, "bluetooth"):
		cmdName = "bluetoothctl"
	case strings.Contains(serviceName, "avahi"):
		cmdName = "avahi-daemon"
	default:
		return "", errors.New("unsupported service")
	}

	// Try to find full path of the command
	fullPath, err := exec.LookPath(cmdName)
	if err != nil {
		// Binary not found in PATH, return empty version but no error
		return "", nil
	}

	var args []string
	switch cmdName {
	case "apache2ctl":
		args = []string{"-v"}
	case "httpd":
		args = []string{"-v"}
	case "nginx":
		args = []string{"-v"}
	case "mysql":
		args = []string{"--version"}
	case "psql":
		args = []string{"--version"}
	case "mongod":
		args = []string{"--version"}
	case "redis-server":
		args = []string{"--version"}
	case "ssh":
		args = []string{"-V"}
	case "docker":
		args = []string{"--version"}
	case "containerd":
		args = []string{"--version"}
	case "crond", "cron":
		args = []string{"--version"}
	case "postfix":
		args = []string{"--version"}
	case "exim":
		args = []string{"-bV"}
	case "dovecot":
		args = []string{"--version"}
	case "memcached":
		args = []string{"-h"} // memcached doesn't have a version flag, -h outputs version info
	case "rabbitmqctl":
		args = []string{"status"} // version is inside status output
	case "firewalld":
		args = []string{"--version"}
	case "ufw":
		args = []string{"version"}
	case "prometheus":
		args = []string{"--version"}
	case "node_exporter":
		args = []string{"--version"}
	case "rsyslogd":
		args = []string{"-v"}
	case "nmcli":
		args = []string{"--version"}
	case "ntpd":
		args = []string{"--version"}
	case "chronyd":
		args = []string{"--version"}
	case "bluetoothctl":
		args = []string{"--version"}
	case "avahi-daemon":
		args = []string{"--version"}
	default:
		// fallback no args
		args = []string{}
	}

	cmd := exec.Command(fullPath, args...)

	outputBytes, err := cmd.CombinedOutput()
	output := string(outputBytes)

	if err != nil {
		// Some commands print version info to stderr (e.g., nginx -v, ssh -V)
		if len(strings.TrimSpace(output)) == 0 {
			return "", err
		}
	}

	// For rabbitmqctl status, version is inside a longer output, extract line containing "RabbitMQ"
	if cmdName == "rabbitmqctl" {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "RabbitMQ") {
				re := regexp.MustCompile(`\d+(\.\d+)+`)
				version := re.FindString(line)
				if version != "" {
					return version, nil
				}
			}
		}
		return "", nil
	}

	re := regexp.MustCompile(`\d+(\.\d+)+`)
	version := re.FindString(output)
	if version == "" {
		return "", nil
	}

	return version, nil
}


func StartService(serviceName string) (string, error) {
	cmd := exec.Command("systemctl", "start", serviceName)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func StopService(serviceName string) (string, error) {
	cmd := exec.Command("systemctl", "stop", serviceName)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func RestartService(serviceName string) (string, error) {
	cmd := exec.Command("systemctl", "restart", serviceName)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func GetServiceLogs(serviceName string) (string, error) {
	cmd := exec.Command("journalctl", "-u", serviceName, "-n", "50", "--no-pager")
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func ServicesActions(conn *websocket.Conn, action, serviceName string) {
	var response struct {
		Type     string `json:"type"`
		Success  bool   `json:"success"`
		Message  string `json:"message"`
		Service  string `json:"service"`
	}

	response.Type = "action_response"
	response.Service = serviceName

	switch action {
	case "start":
		output, err := StartService(serviceName)
		if err != nil {
			response.Success = false
			response.Message = "Failed to start service: " + err.Error()
			log.Printf("Error starting service %s: %v", serviceName, err)
		} else {
			response.Success = true
			response.Message = "Service started successfully: " + output
		}
	case "stop":
		output, err := StopService(serviceName)
		if err != nil {
			response.Success = false
			response.Message = "Failed to stop service: " + err.Error()
			log.Printf("Error stopping service %s: %v", serviceName, err)
		} else {
			response.Success = true
			response.Message = "Service stopped successfully: " + output
		}
	case "restart":
		output, err := RestartService(serviceName)
		if err != nil {
			response.Success = false
			response.Message = "Failed to restart service: " + err.Error()
			log.Printf("Error restarting service %s: %v", serviceName, err)
		} else {
			response.Success = true
			response.Message = "Service restarted successfully: " + output
		}
	case "logs":
		logs, err := GetServiceLogs(serviceName)
		if err != nil {
			response.Success = false
			response.Message = "Error getting logs: " + err.Error()
			log.Printf("Error getting logs for service %s: %v", serviceName, err)
		} else {
			logResponse := struct {
				Type    string `json:"type"`
				Service string `json:"service"`
				Log     string `json:"log"`
			}{
				Type:    "log",
				Service: serviceName,
				Log:     logs,
			}
			data, _ := json.Marshal(logResponse)
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Error sending log response: %v", err)
			}
			return
		}
	default:
		response.Success = false
		response.Message = "Invalid service action: " + action
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}