package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"regexp"
	"github.com/gorilla/websocket"
)

type FirewallRule struct {
	ID          string `json:"id"`
	Protocol    string `json:"protocol"`
	Port        string `json:"port"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Action      string `json:"action"`
}

type RunningPort struct {
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
	Process  string `json:"process"`
	PID      int    `json:"pid"`
}

// Parses `ss -tulnp` output to extract protocol, port, process, and PID
func parseSSOutput(output string) []RunningPort {
	lines := strings.Split(output, "\n")
	var ports []RunningPort

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "tcp") && !strings.HasPrefix(line, "udp") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}

		// Extract protocol
		proto := fields[0]

		// Extract port from local address (field 4)
		localAddr := fields[4]
		parts := strings.Split(localAddr, ":")
		portStr := parts[len(parts)-1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Printf("Failed to parse port %s: %v", portStr, err)
			continue
		}

		// Extract process and PID from the last field (usually contains "users:(...)"")
		var process string
		var pid int
		lastField := fields[len(fields)-1]
		if strings.HasPrefix(lastField, "users:((") {
			// Remove "users:((" prefix and "))" suffix
			userInfo := strings.TrimPrefix(lastField, "users:((")
			userInfo = strings.TrimSuffix(userInfo, "))")
			// Split by comma to get process and PID
			userParts := strings.Split(userInfo, ",")
			if len(userParts) >= 2 {
				// Process name is usually the first part
				process = strings.Split(userParts[0], "\"")[1]
				// PID is usually in the second part, after "pid="
				for _, part := range userParts {
					if strings.HasPrefix(part, "pid=") {
						pidStr := strings.TrimPrefix(part, "pid=")
						pid, err = strconv.Atoi(pidStr)
						if err != nil {
							log.Printf("Failed to parse PID %s: %v", pidStr, err)
							continue
						}
						break
					}
				}
			}
		}

		ports = append(ports, RunningPort{
			Protocol: proto,
			Port:     port,
			Process:  process,
			PID:      pid,
		})
	}
	return ports
}

// Checks if a process is still running
func isProcessRunning(pid int) bool {
	if pid == 0 {
		return false
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	if runtime.GOOS == "windows" {
		// On Windows, FindProcess always succeeds, so we need to check if the process is actually running
		// Attempt to open the process with minimal permissions
		// This is a simplistic check; consider using Windows API for more accuracy if needed
		cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid))
		output, err := cmd.CombinedOutput()
		if err != nil {
			return false
		}
		return strings.Contains(string(output), fmt.Sprintf("%d", pid))
	}
	// On Unix-like systems, Signal(0) checks if the process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

var lastErrorTime time.Time
var errorLogInterval = 1 * time.Minute

func logErrorWithRateLimit(message string, err error) {
    if time.Since(lastErrorTime) > errorLogInterval {
        log.Printf("%s: %v", message, err)
        lastErrorTime = time.Now()
    }
}

func GetAllFirewallRules() []FirewallRule {
    cmd := exec.Command("sudo", "ufw", "status", "numbered")
    output, err := cmd.CombinedOutput()
    if err != nil {
        logErrorWithRateLimit("Error fetching ufw firewall rules", err)
        return []FirewallRule{}
    }

    lines := strings.Split(string(output), "\n")
    var rules []FirewallRule

    // Sample ufw numbered output line example:
    // [ 1] 22/tcp                     ALLOW IN    Anywhere
    // [ 2] 80/tcp                     DENY IN     192.168.0.0/24

    // Skip header lines until you find rules (usually after "Status: active" line)
    startParsing := false
    id := 1

    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }
        if strings.HasPrefix(line, "Status:") {
            // Status line, just note active or inactive
            if !strings.Contains(line, "active") {
                // Firewall inactive, return empty
                return []FirewallRule{}
            }
            startParsing = true
            continue
        }
        if !startParsing {
            continue
        }

        // Rule lines start with: [ <num>] ...
        if !strings.HasPrefix(line, "[") {
            continue
        }

        // Parse with regexp to extract id, port/protocol, action, source
        // Example line: [ 1] 22/tcp                     ALLOW IN    Anywhere
        re := regexp.MustCompile(`\[\s*(\d+)\]\s+(\S+)\s+(\S+)\s+(\S+)\s+(.*)`)
        matches := re.FindStringSubmatch(line)
        if len(matches) < 6 {
            // Try simpler fallback: split by spaces
            fields := strings.Fields(line)
            if len(fields) < 4 {
                continue
            }
            // Fallback parsing - crude
            portProto := fields[1]
            action := fields[2]
            source := fields[len(fields)-1]
            protocol := "unknown"
            port := ""

            if strings.Contains(portProto, "/") {
                parts := strings.Split(portProto, "/")
                port = parts[0]
                protocol = parts[1]
            } else {
                port = portProto
            }

            rule := FirewallRule{
                ID:          strconv.Itoa(id),
                Action:      action,
                Protocol:    protocol,
                Port:        port,
                Source:      source,
                Destination: "any",
            }
            rules = append(rules, rule)
            id++
            continue
        }

        // Normal parse
        idStr := matches[1]
        portProto := matches[2]
        action := matches[3]
        // matches[4] is usually IN or OUT direction, can ignore or keep if needed
        source := matches[5]

        protocol := "unknown"
        port := ""
        if strings.Contains(portProto, "/") {
            parts := strings.Split(portProto, "/")
            port = parts[0]
            protocol = parts[1]
        } else {
            port = portProto
        }

        rule := FirewallRule{
            ID:          idStr,
            Action:      action,
            Protocol:    strings.ToLower(protocol),
            Port:        port,
            Source:      source,
            Destination: "any",
        }
        rules = append(rules, rule)
        id++
    }
    return rules
}

func GetAllRunningPorts() ([]RunningPort, error) {
	// Use `ss -tulnp` to include process and PID information
	cmd := exec.Command("ss", "-tulnp")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing ss -tulnp: %v, output: %s", err, string(output))
		return nil, err
	}
	return parseSSOutput(string(output)), nil
}

func StopPort(port int, pid int) error {
	if pid == 0 {
		return fmt.Errorf("invalid PID: %d", pid)
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid port: %d", port)
	}

	log.Printf("Attempting to stop process %d on port %d", pid, port)

	// Verify the process is still running and associated with the port
	if !isProcessRunning(pid) {
		return fmt.Errorf("process %d is not running", pid)
	}

	// Verify the port is still in use by the process
	ports, err := GetAllRunningPorts()
	if err != nil {
		return fmt.Errorf("failed to verify port usage: %v", err)
	}
	portInUse := false
	for _, p := range ports {
		if p.PID == pid && p.Port == port {
			portInUse = true
			break
		}
	}
	if !portInUse {
		return fmt.Errorf("port %d is not in use by PID %d", port, pid)
	}

	// Find the process
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %v", pid, err)
	}

	// Send SIGTERM and wait
	if err := process.Signal(syscall.SIGTERM); err != nil {
		log.Printf("Failed to send SIGTERM to PID %d: %v", pid, err)
		// Attempt SIGKILL as a fallback
		if err := process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process %d: %v", pid, err)
		}
		return nil
	}

	// Wait for the process to terminate (up to 5 seconds)
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			log.Printf("Process %d did not terminate within 5 seconds, sending SIGKILL", pid)
			if err := process.Kill(); err != nil {
				return fmt.Errorf("failed to kill process %d: %v", pid, err)
			}
			return nil
		case <-ticker.C:
			if !isProcessRunning(pid) {
				log.Printf("Process %d terminated successfully", pid)
				return nil
			}
		}
	}
}

func AddPort(port int, protocol string) error {
	log.Printf("Attempting to add port %d (%s)", port, protocol)

	protocol = strings.ToUpper(protocol)
	if protocol != "TCP" && protocol != "UDP" {
		return fmt.Errorf("invalid protocol: %s", protocol)
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid port: %d", port)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		ruleName := fmt.Sprintf("Grok-Open-Port-%d-%s", port, protocol)
		cmd = exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
			"name="+ruleName,
			"dir=in", "action=allow", "protocol="+protocol, "localport="+strconv.Itoa(port))
	} else {
		cmd = exec.Command("ufw", "allow", fmt.Sprintf("%d/%s", port, strings.ToLower(protocol)))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error adding port: %v, output: %s", err, string(output))
		return err
	}
	return nil
}

func MonitorAndBroadcastFirewallRules(broadcast chan<- []byte) {
    for {
        rules := GetAllFirewallRules()
        data, err := json.Marshal(struct {
            Type  string         `json:"type"`
            Rules []FirewallRule `json:"rules"`
        }{
            Type:  "firewalls",
            Rules: rules,
        })
        if err != nil {
            logErrorWithRateLimit("Error marshaling firewall rules", err)
            time.Sleep(10 * time.Second)
            continue
        }
        select {
        case broadcast <- data:
            log.Println("Sent firewall rules to broadcast channel")
        default:
            log.Println("Broadcast channel full, skipping firewall rules update")
        }
        time.Sleep(30 * time.Second) // Poll every 30 seconds
    }
}

func MonitorAndBroadcastRunningPorts(broadcast chan<- []byte) {
	for {
		ports, err := GetAllRunningPorts()
		if err != nil {
			log.Printf("Error fetching running ports: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		data, err := json.Marshal(struct {
			Type  string        `json:"type"`
			Ports []RunningPort `json:"ports"`
		}{
			Type:  "ports",
			Ports: ports,
		})
		if err != nil {
			log.Printf("Error marshaling running ports: %v", err)
			continue
		}
		broadcast <- data
		time.Sleep(5 * time.Second)
	}
}

func FirewallActions(conn *websocket.Conn, actionType string, data json.RawMessage) {
	var response struct {
		Type    string `json:"type"`
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	response.Type = "action_response"

	switch actionType {
	case "stop_port":
		var portData struct {
			Port int `json:"port"`
			PID  int `json:"pid"`
		}
		if err := json.Unmarshal(data, &portData); err == nil {
			if err := StopPort(portData.Port, portData.PID); err == nil {
				response.Success = true
				response.Message = fmt.Sprintf("Port %d stopped successfully", portData.Port)
			} else {
				response.Success = false
				response.Message = fmt.Sprintf("Failed to stop port %d: %v", portData.Port, err)
			}
		} else {
			response.Success = false
			response.Message = "Invalid port data: " + err.Error()
		}

	case "add_port":
		var portData struct {
			Port     int    `json:"port"`
			Protocol string `json:"protocol"`
		}
		if err := json.Unmarshal(data, &portData); err == nil {
			if err := AddPort(portData.Port, portData.Protocol); err == nil {
				response.Success = true
				response.Message = fmt.Sprintf("Port %d/%s added successfully", portData.Port, portData.Protocol)
			} else {
				response.Success = false
				response.Message = fmt.Sprintf("Failed to add port %d/%s: %v", portData.Port, portData.Protocol, err)
			}
		} else {
			response.Success = false
			response.Message = "Invalid port data: " + err.Error()
		}

	case "firewall_add_rule":
		// Example: Add a new firewall rule (future implementation)
		var ruleData struct {
			Protocol    string `json:"protocol"`
			Port        string `json:"port"`
			Source      string `json:"source"`
			Destination string `json:"destination"`
			Action      string `json:"action"`
		}
		if err := json.Unmarshal(data, &ruleData); err == nil {
			if err := AddFirewallRule(ruleData.Protocol, ruleData.Port, ruleData.Source, ruleData.Destination, ruleData.Action); err == nil {
				response.Success = true
				response.Message = fmt.Sprintf("Firewall rule for %s:%s added successfully", ruleData.Protocol, ruleData.Port)
			} else {
				response.Success = false
				response.Message = fmt.Sprintf("Failed to add firewall rule: %v", err)
			}
		} else {
			response.Success = false
			response.Message = "Invalid firewall rule data: " + err.Error()
		}

	default:
		response.Success = false
		response.Message = "Unknown firewall action: " + actionType
	}

	dataResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, dataResponse); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func AddFirewallRule(protocol, port, source, destination, action string) error {
	if protocol != "tcp" && protocol != "udp" {
		return fmt.Errorf("invalid protocol: %s", protocol)
	}
	if port == "" {
		return fmt.Errorf("port cannot be empty")
	}
	if action != "allow" && action != "deny" {
		return fmt.Errorf("invalid action: %s", action)
	}

	cmdStr := fmt.Sprintf("%s/%s", port, strings.ToLower(protocol))
	if source != "" && source != "any" {
		cmdStr = fmt.Sprintf("%s from %s", cmdStr, source)
	}
	if destination != "" && destination != "any" {
		cmdStr = fmt.Sprintf("%s to %s", cmdStr, destination)
	}

	cmd := exec.Command("sudo", "ufw", strings.ToLower(action), cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error adding firewall rule: %v, output: %s", err, string(output))
		return err
	}
	return nil
}