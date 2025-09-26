package routes

import (
    "net/http"
    "time"

    "snaptrack/db"

    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

func RegisterDashboardRoutes(app *fiber.App) {
    app.Get("/api/dashboard/stats", getDashboardStats)
    app.Get("/api/dashboard/recent-activity", getRecentActivity)
}

// Dashboard stats response
type DashboardStatsResponse struct {
    TotalBackups int64      `json:"total_backups"`
    StorageUsed  int64      `json:"storage_used"`
    LastBackup   *time.Time `json:"last_backup"`
    SystemStatus string     `json:"system_status"`
}

func getDashboardStats(c *fiber.Ctx) error {
    var totalBackups int64
    var storageUsed int64
    var lastBackup db.Backup

    // Count all backups
    if err := db.DB.Model(&db.Backup{}).Count(&totalBackups).Error; err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count backups"})
    }

    // Sum storage used
    if err := db.DB.Model(&db.Backup{}).Select("COALESCE(SUM(size_bytes), 0)").Scan(&storageUsed).Error; err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to calculate storage"})
    }

    // Get last backup (most recent completed)
    if err := db.DB.Where("completed_at IS NOT NULL").Order("completed_at desc").First(&lastBackup).Error; err != nil {
        if err != gorm.ErrRecordNotFound {
            return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch last backup"})
        }
    }

    // Determine system status (online if at least one enabled server)
    var enabledServers int64
    if err := db.DB.Model(&db.Server{}).Where("enabled = ?", true).Count(&enabledServers).Error; err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check servers"})
    }
    systemStatus := "offline"
    if enabledServers > 0 {
        systemStatus = "online"
    }

    return c.JSON(DashboardStatsResponse{
        TotalBackups: totalBackups,
        StorageUsed:  storageUsed,
        LastBackup:   lastBackup.CompletedAt,
        SystemStatus: systemStatus,
    })
}

func getRecentActivity(c *fiber.Ctx) error {
    var logs []db.Log
    if err := db.DB.Order("created_at desc").Limit(10).Find(&logs).Error; err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch activity"})
    }

    return c.JSON(fiber.Map{
        "activities": logs,
    })
}
