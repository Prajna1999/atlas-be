package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Prajna1999/atlas-be/internal/database"
	"github.com/Prajna1999/atlas-be/internal/routes"
	"github.com/Prajna1999/atlas-be/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type App struct {
	router   *gin.Engine
	db       *gorm.DB
	routes   *routes.Routes
	services map[string]interface{}
}

// NewApp initializes and returns a pointer to the App.
func NewApp() (*App, error) {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize DB: %v", err)
	}

	// Run database migrations
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate DB: %v", err)
	}

	// Initialize routes and services
	services, err := initializeServices(db)
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	hetznerService := services["hetzner"].(*service.HetznerService)
	routes := routes.NewRoutes(hetznerService)

	// Create the App instance
	app := &App{
		router:   gin.Default(),
		db:       db,
		routes:   routes,
		services: services,
	}

	app.setupRoutes()
	return app, nil
}

func (a *App) setupRoutes() {
	// Health check route
	a.router.GET("/api/v1/health-check", a.healthCheck)

	// Set up application routes
	a.routes.SetupRoutes(a.router)
}

// Run starts the application server on the default port.
func (a *App) Run() error {
	return a.router.Run(":8080")
}

// healthCheck verifies the application's health, including the database connection.
func (a *App) healthCheck(c *gin.Context) {
	response := HealthCheckResponse{
		Status:  "OK",
		Message: "All Systems Normal 🚀",
	}

	sqlDB, err := a.db.DB()
	if err != nil {
		response.Status = "Error"
		response.Message = "Database Connection Error"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Ping the database to verify connection
	if err := sqlDB.Ping(); err != nil {
		response.Status = "Error"
		response.Message = "Database Ping Failed"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// initializeServices initializes and returns all services in a map.
func initializeServices(db *gorm.DB) (map[string]interface{}, error) {
	services := make(map[string]interface{})

	// Initialize HetznerService
	hetznerToken := os.Getenv("HCLOUD_TOKEN")
	// fmt.Printf("the hcloud toke is %v", hetznerToken)
	if hetznerToken == "" {
		return nil, fmt.Errorf("HCLOUD_TOKEN environment variable not set")
	}

	hetznerService, err := service.NewHetznerService(hetznerToken)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Hetzner service: %v", err)
	}

	services["hetzner"] = hetznerService

	// Add other service initializations here

	return services, nil
}
