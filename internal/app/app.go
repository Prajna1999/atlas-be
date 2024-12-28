package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Prajna1999/atlas-be/internal/database"
	"github.com/Prajna1999/atlas-be/internal/routes"
	"github.com/Prajna1999/atlas-be/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type HealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type App struct {
	router   *gin.Engine
	db       *database.DBClient
	routes   *routes.Routes
	services map[string]interface{}
}

// NewApp initializes and returns a pointer to the App.
func NewApp() (*App, error) {
	// Initialize MongoDB client
	db, err := database.InitDB()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize collections: %v", err)
	}

	// Initialize the collections
	if err := initializeCollections(db); err != nil {
		return nil, fmt.Errorf("failed to initialize collections: %v", err)
	}

	// Initialize routes and services
	services, err := initializeServices(db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize services: %v", err)
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

// initialise collections
func initializeCollections(db *database.DBClient) error {
	// Get the Atlas database
	atlasDB := db.GetDatabase("atlas-test")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// List existing collections
	collections, err := atlasDB.ListCollectionNames(ctx, bson.D{})
	fmt.Printf("the connections are %v", collections)
	if err != nil {
		return fmt.Errorf("failed to list collections: %v", err)
	}

	// Create collections if they don't exist
	requiredCollections := []string{"users"}

	for _, collName := range requiredCollections {
		// Check if collection already exists
		exists := false
		for _, existingColl := range collections {
			if existingColl == collName {
				exists = true
				break
			}
		}

		// Skip if collection already exists
		if exists {
			log.Printf("Collection %s already exists", collName)
			continue
		}

		// Create collection
		command := bson.D{{Key: "create", Value: collName}}
		err := atlasDB.RunCommand(ctx, command).Err()
		if err != nil {
			return fmt.Errorf("failed to create collection %s: %v", collName, err)
		}

		log.Printf("Created collection: %s", collName)
	}

	return nil
}

// // isCollectionExistsError checks if the error is due to collection already existing
//
//	func isCollectionExistsError(err error) bool {
//		return err != nil && err.Error() == "namespace already exists"
//	}
func (a *App) setupRoutes() {
	// Health check route
	a.router.GET("/api/v1/health-check", a.healthCheck)

	// Set up application routes
	a.routes.SetupRoutes(a.router)
}

// Run starts the application server on the default port.
func (a *App) Run() error {

	defer func() {
		if err := a.db.Close(); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}()

	return a.router.Run(":8080")
}

// healthCheck verifies the application's health, including the database connection.
func (a *App) healthCheck(c *gin.Context) {
	response := HealthCheckResponse{
		Status:  "OK",
		Message: "All Systems Normal 🚀",
	}

	// create a context with timeout for the health check
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping MongoDB to verify connection
	if err := a.db.Client.Ping(ctx, nil); err != nil {
		response.Status = "Error"
		response.Message = "Database Connection Error"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// initializeServices initializes and returns all services in a map.
func initializeServices(db *database.DBClient) (map[string]interface{}, error) {
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

	// Initialize MongoDB-specific services
	//  userService := service.NewUserService(db.GetCollection("atlas", "users"))
	//  services["user"] = userService

	//  projectService := service.NewProjectService(db.GetCollection("atlas", "projects"))
	//  services["project"] = projectService

	//  serverService := service.NewServerService(db.GetCollection("atlas", "servers"))
	//  services["server"] = serverService

	return services, nil
}
