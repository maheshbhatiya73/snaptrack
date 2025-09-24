package routes

import (
	"fmt"
	"net"
	"os"
	"snaptrack/auth"
	"snaptrack/db"
	"strings"
	"time"

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

// -------------------- Helper Functions --------------------

func createSSHClient(host string, user, keyPath *string, port *int) (*ssh.Client, error) {
	if user == nil || keyPath == nil || port == nil {
		return nil, fmt.Errorf("SSH configuration incomplete")
	}

	key, err := os.ReadFile(*keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH key: %v", err)
	}

	config := &ssh.ClientConfig{
		User:            *user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, *port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("SSH connection failed: %v", err)
	}

	return client, nil
}

func validateRemoteServer(host string, sshUser, sshKeyPath *string, sshPort *int, transferType *string) error {
	// TCP reachability check
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, "80"), 5*time.Second)
	if err != nil {
		conn, err = net.DialTimeout("tcp", net.JoinHostPort(host, "443"), 5*time.Second)
		if err != nil {
			return fmt.Errorf("host %s is not reachable: %v", host, err)
		}
	}
	if conn != nil {
		conn.Close()
	}

	client, err := createSSHClient(host, sshUser, sshKeyPath, sshPort)
	if err != nil {
		return err
	}
	defer client.Close()

	if transferType != nil && *transferType != "rsync" && *transferType != "scp" {
		return fmt.Errorf("unsupported transfer type: %s", *transferType)
	}

	return nil
}

// -------------------- Route Handlers --------------------

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

func createServer(c *fiber.Ctx) error {
	var server db.Server
	if err := c.BodyParser(&server); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var existingServer db.Server
	if err := db.DB.Unscoped().Where("name = ?", server.Name).First(&existingServer).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Server name already exists"})
	}

	if server.Type == "remote" {
		if err := validateRemoteServer(server.Host, server.SSHUser, server.SSHKeyPath, server.SSHPort, server.TransferType); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Server validation failed: %v", err)})
		}
	}

	if err := db.DB.Create(&server).Error; err != nil {
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

	// Check for name conflict
	if updateData.Name != "" && updateData.Name != server.Name {
		var existingServer db.Server
		if err := db.DB.Unscoped().Where("name = ? AND id != ?", updateData.Name, server.ID).First(&existingServer).Error; err == nil {
			return c.Status(409).JSON(fiber.Map{"error": "Server name already exists"})
		}
	}

	// Resolve final values
	finalType := server.Type
	if updateData.Type != "" {
		finalType = updateData.Type
	}
	finalHost := server.Host
	if updateData.Host != "" {
		finalHost = updateData.Host
	}
	finalSSHUser := server.SSHUser
	if updateData.SSHUser != nil {
		finalSSHUser = updateData.SSHUser
	}
	finalSSHPort := server.SSHPort
	if updateData.SSHPort != nil {
		finalSSHPort = updateData.SSHPort
	}
	finalSSHKeyPath := server.SSHKeyPath
	if updateData.SSHKeyPath != nil {
		finalSSHKeyPath = updateData.SSHKeyPath
	}
	finalTransferType := server.TransferType
	if updateData.TransferType != nil {
		finalTransferType = updateData.TransferType
	}

	if finalType == "remote" {
		if err := validateRemoteServer(finalHost, finalSSHUser, finalSSHKeyPath, finalSSHPort, finalTransferType); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("Server validation failed: %v", err)})
		}
	} else {
		localType := "local"
		finalTransferType = &localType
	}

	updateData.TransferType = finalTransferType

	if err := db.DB.Model(&server).Updates(updateData).Error; err != nil {
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

	client, err := createSSHClient(server.Host, server.SSHUser, server.SSHKeyPath, server.SSHPort)
	if err != nil {
		return c.JSON(fiber.Map{"success": false, "message": err.Error()})
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
        if _, err := os.Stat(req.Path); os.IsNotExist(err) {
            return c.JSON(fiber.Map{"valid": false, "message": "Path does not exist"})
        }
        return c.JSON(fiber.Map{"valid": true, "message": "Path exists"})
    }

    client, err := createSSHClient(server.Host, server.SSHUser, server.SSHKeyPath, server.SSHPort)
    if err != nil {
        return c.JSON(fiber.Map{"valid": false, "message": err.Error()})
    }
    defer client.Close()

    session, err := client.NewSession()
    if err != nil {
        return c.JSON(fiber.Map{"valid": false, "message": fmt.Sprintf("Failed to create SSH session: %v", err)})
    }
    defer session.Close()

    escapedPath := strings.ReplaceAll(req.Path, `"`, `\"`)
    cmd := fmt.Sprintf(`test -d "%s" && test -r "%s" && test -x "%s"`, escapedPath, escapedPath, escapedPath)

    done := make(chan error, 1)
    go func() {
        _, err := session.CombinedOutput(cmd)
        done <- err
    }()

    select {
    case err := <-done:
        if err != nil {
            return c.JSON(fiber.Map{"valid": false, "message": fmt.Sprintf("Path '%s' is not accessible: %v", req.Path, err)})
        }
    case <-time.After(5 * time.Second):
        return c.JSON(fiber.Map{"valid": false, "message": "Path check timed out"})
    }

    return c.JSON(fiber.Map{"valid": true, "message": "Path exists and is accessible"})
}

// ValidateLocalPath validates a local filesystem path on the server running this backend
func ValidateLocalPath(c *fiber.Ctx) error {
    var req struct {
        Path string `json:"path"`
    }
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }

    if strings.TrimSpace(req.Path) == "" {
        return c.Status(400).JSON(fiber.Map{"error": "path is required"})
    }

    if _, err := os.Stat(req.Path); os.IsNotExist(err) {
        return c.JSON(fiber.Map{"valid": false, "message": "Path does not exist"})
    } else if err != nil {
        return c.JSON(fiber.Map{"valid": false, "message": fmt.Sprintf("Failed to access path: %v", err)})
    }

    return c.JSON(fiber.Map{"valid": true, "message": "Path exists"})
}
