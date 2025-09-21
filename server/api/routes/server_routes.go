package routes

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	"snaptrack/auth"
	"snaptrack/db"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/ssh"
)

func RegisterServerRoutes(app *fiber.App) {
	api := app.Group("/api/servers", auth.RequireJWT())

	api.Get("/", listServers)
	api.Get("/:id", getServer)
	api.Post("/", createServer)
	api.Put("/:id", updateServer)
	api.Delete("/:id", deleteServer)
	api.Post("/:id/test", testServerConnection)
	api.Post("/:id/validate-path", validatePath)
}

func listServers(c *fiber.Ctx) error {
	var servers []db.Server
	db.DB.Find(&servers)
	return c.JSON(servers)
}

func getServer(c *fiber.Ctx) error {
	id := c.Params("id")
	var server db.Server
	if err := db.DB.First(&server, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Server not found"})
	}
	return c.JSON(server)
}

// validateRemoteServer performs validation for remote servers
func validateRemoteServer(host string, sshUser, sshKeyPath *string, sshPort *int) error {
	// First, ping the host to check if it's reachable
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "80"), 5*time.Second)
	if err != nil {
		// Try common ports if 80 fails
		conn, err = net.DialTimeout("tcp", net.JoinHostPort(host, "443"), 5*time.Second)
		if err != nil {
			return fmt.Errorf("host %s is not reachable: %v", host, err)
		}
	}
	if conn != nil {
		conn.Close()
	}

	// Validate SSH connection
	if sshUser == nil || sshKeyPath == nil || sshPort == nil {
		return fmt.Errorf("SSH configuration incomplete")
	}

	key, err := os.ReadFile(*sshKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read SSH key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to parse SSH key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: *sshUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, *sshPort)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("SSH connection failed: %v", err)
	}
	defer client.Close()

	return nil
}

func createServer(c *fiber.Ctx) error {
	var server db.Server
	if err := c.BodyParser(&server); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if server name already exists (including soft-deleted records)
	var existingServer db.Server
	if err := db.DB.Unscoped().Where("name = ?", server.Name).First(&existingServer).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Server name already exists"})
	}

	// Validate remote server if type is remote
	if server.Type == "remote" {
		if err := validateRemoteServer(server.Host, server.SSHUser, server.SSHKeyPath, server.SSHPort); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Server validation failed: %v", err)})
		}
	}

	if err := db.DB.Create(&server).Error; err != nil {
		// Check if it's a duplicate key error
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return c.Status(409).JSON(fiber.Map{"error": "Server name already exists"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(server)
}

func updateServer(c *fiber.Ctx) error {
	id := c.Params("id")
	var server db.Server
	if err := db.DB.First(&server, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Server not found"})
	}

	var updateData db.Server
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if the new name already exists (excluding current server and including soft-deleted)
	if updateData.Name != "" && updateData.Name != server.Name {
		var existingServer db.Server
		if err := db.DB.Unscoped().Where("name = ? AND id != ?", updateData.Name, server.ID).First(&existingServer).Error; err == nil {
			return c.Status(409).JSON(fiber.Map{"error": "Server name already exists"})
		}
	}

	// Validate remote server configuration if updating to remote or updating remote fields
	finalType := updateData.Type
	if finalType == "" {
		finalType = server.Type
	}

	finalHost := updateData.Host
	if finalHost == "" {
		finalHost = server.Host
	}

	finalSSHUser := updateData.SSHUser
	if finalSSHUser == nil {
		finalSSHUser = server.SSHUser
	}

	finalSSHPort := updateData.SSHPort
	if finalSSHPort == nil {
		finalSSHPort = server.SSHPort
	}

	finalSSHKeyPath := updateData.SSHKeyPath
	if finalSSHKeyPath == nil {
		finalSSHKeyPath = server.SSHKeyPath
	}

	if finalType == "remote" {
		if err := validateRemoteServer(finalHost, finalSSHUser, finalSSHKeyPath, finalSSHPort); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Server validation failed: %v", err)})
		}
	}

	if err := db.DB.Model(&server).Updates(updateData).Error; err != nil {
		// Check if it's a duplicate key error
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return c.Status(409).JSON(fiber.Map{"error": "Server name already exists"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(server)
}

func deleteServer(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.DB.Delete(&db.Server{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}

func testServerConnection(c *fiber.Ctx) error {
	id := c.Params("id")
	var server db.Server
	if err := db.DB.First(&server, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Server not found"})
	}

	if server.Type == "local" {
		return c.JSON(fiber.Map{"success": true, "message": "Local server connection is always available"})
	}

	// Test SSH connection for remote servers
	if server.SSHUser == nil || server.SSHKeyPath == nil || server.SSHPort == nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "SSH configuration incomplete"})
	}

	key, err := os.ReadFile(*server.SSHKeyPath)
	if err != nil {
		return c.JSON(fiber.Map{"success": false, "message": fmt.Sprintf("Failed to read SSH key: %v", err)})
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return c.JSON(fiber.Map{"success": false, "message": fmt.Sprintf("Failed to parse SSH key: %v", err)})
	}

	config := &ssh.ClientConfig{
		User: *server.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", server.Host, *server.SSHPort)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return c.JSON(fiber.Map{"success": false, "message": fmt.Sprintf("SSH connection failed: %v", err)})
	}
	defer client.Close()

	return c.JSON(fiber.Map{"success": true, "message": "SSH connection successful"})
}

func validatePath(c *fiber.Ctx) error {
	id := c.Params("id")
	var server db.Server
	if err := db.DB.First(&server, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Server not found"})
	}

	var req struct {
		Path string `json:"path"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if server.Type == "local" {
		// Check if path exists on local filesystem
		if _, err := os.Stat(req.Path); os.IsNotExist(err) {
			return c.JSON(fiber.Map{"valid": false, "message": "Path does not exist"})
		}
		return c.JSON(fiber.Map{"valid": true, "message": "Path exists"})
	}

	// Validate path on remote server via SSH
	if server.SSHUser == nil || server.SSHKeyPath == nil || server.SSHPort == nil {
		return c.Status(400).JSON(fiber.Map{"valid": false, "message": "SSH configuration incomplete"})
	}

	key, err := os.ReadFile(*server.SSHKeyPath)
	if err != nil {
		return c.JSON(fiber.Map{"valid": false, "message": fmt.Sprintf("Failed to read SSH key: %v", err)})
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return c.JSON(fiber.Map{"valid": false, "message": fmt.Sprintf("Failed to parse SSH key: %v", err)})
	}

	config := &ssh.ClientConfig{
		User: *server.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", server.Host, *server.SSHPort)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return c.JSON(fiber.Map{"valid": false, "message": fmt.Sprintf("SSH connection failed: %v", err)})
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return c.JSON(fiber.Map{"valid": false, "message": fmt.Sprintf("Failed to create SSH session: %v", err)})
	}
	defer session.Close()

	// Check if path exists using ls command
	cmd := fmt.Sprintf("ls -ld '%s'", req.Path)
	_, err = session.CombinedOutput(cmd)
	if err != nil {
		return c.JSON(fiber.Map{"valid": false, "message": fmt.Sprintf("Path does not exist or is not accessible: %v", err)})
	}

	return c.JSON(fiber.Map{"valid": true, "message": "Path exists and is accessible"})
}
