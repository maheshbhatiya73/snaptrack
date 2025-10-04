package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"snaptrack/api"
	"snaptrack/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config struct
type Config struct {
	Env    string `yaml:"env"`
	Server struct {
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		FrontendPath string `yaml:"frontend_path"`
	} `yaml:"server"`
	CORS struct {
		Origins string `yaml:"origins"`
	} `yaml:"cors"`
	Security struct {
		JWTSecret string `yaml:"jwt_secret"`
	} `yaml:"security"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
}

// LoadConfig loads configuration from /etc/snaptrack/config.yaml or fallback
func LoadConfig() *Config {
	config := &Config{}

	if _, err := os.Stat("/etc/snaptrack/config.yaml"); err == nil {
		data, err := ioutil.ReadFile("/etc/snaptrack/config.yaml")
		if err != nil {
			log.Fatalf("Failed to read config.yaml: %v", err)
		}
		if err := yaml.Unmarshal(data, config); err != nil {
			log.Fatalf("Failed to parse config.yaml: %v", err)
		}
		log.Println("Loaded config from /etc/snaptrack/config.yaml")
		return config
	}

	_ = godotenv.Load() // fallback .env

	// Defaults / env
	config.Env = os.Getenv("ENV")
	if config.Env == "" {
		config.Env = "development"
	}

	config.Server.Host = os.Getenv("HOST")
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}

	config.Server.Port = os.Getenv("PORT")
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}

	config.Server.FrontendPath = os.Getenv("FRONTEND_PATH")
	if config.Server.FrontendPath == "" {
		config.Server.FrontendPath = "./web/.output/public"
	}

	config.CORS.Origins = os.Getenv("CORS_ORIGINS")
	if config.CORS.Origins == "" {
		config.CORS.Origins = "*"
	}

	config.Security.JWTSecret = os.Getenv("JWT_SECRET")
	if config.Security.JWTSecret == "" {
		config.Security.JWTSecret = "snaptrack"
	}

	config.Database.Host = os.Getenv("PG_HOST")
	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}
	config.Database.Port = 5432
	config.Database.User = os.Getenv("PG_USER")
	if config.Database.User == "" {
		config.Database.User = "postgres"
	}
	config.Database.Password = os.Getenv("PG_PASSWORD")
	if config.Database.Password == "" {
		config.Database.Password = "mahesh"
	}
	config.Database.DBName = os.Getenv("PG_DBNAME")
	if config.Database.DBName == "" {
		config.Database.DBName = "snaptrack"
	}

	log.Println("Loaded config from environment variables / defaults")
	return config
}

func main() {
	config := LoadConfig()

	// Set env vars for DB
	os.Setenv("PG_HOST", config.Database.Host)
	os.Setenv("PG_PORT", fmt.Sprintf("%d", config.Database.Port))
	os.Setenv("PG_USER", config.Database.User)
	os.Setenv("PG_PASSWORD", config.Database.Password)
	os.Setenv("PG_DBNAME", config.Database.DBName)

	// Connect to DB
	db.Connect()

	// Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:          10 * time.Second,
		WriteTimeout:         10 * time.Second,
		DisableStartupMessage: true, // hide banner
	})

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.CORS.Origins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// === API routes first ===
	api.RegisterRoutes(app)
	api.RegisterWebSocketRoutes(app)

	// === Serve frontend SPA ===
	if _, err := os.Stat(config.Server.FrontendPath); os.IsNotExist(err) {
		log.Printf("⚠ Frontend folder not found: %s\n", config.Server.FrontendPath)
	} else {
		app.Static("/", config.Server.FrontendPath)
		app.Use(func(c *fiber.Ctx) error {
			if c.Path()[:5] != "/api/" { // SPA fallback for non-API
				return c.SendFile(fmt.Sprintf("%s/index.html", config.Server.FrontendPath))
			}
			return c.Next()
		})
		log.Printf("✅ Frontend served from: %s\n", config.Server.FrontendPath)
	}

	addr := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
	log.Printf("Server running at http://%s\n", addr)
	log.Fatal(app.Listen(addr))
}
