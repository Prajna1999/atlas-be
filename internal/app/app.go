package app

import (
	"log"
	"net/http"

	"github.com/Prajna1999/hetzner-be/internal/database"
	"github.com/Prajna1999/hetzner-be/internal/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
type App struct {
	router *gin.Engine
	db     *gorm.DB

	routes *routes.Routes
}

// function that returns a pointer to the new app
func NewApp() (*App, error) {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize DB %v", err)
	}

	err = db.AutoMigrate()
	if err != nil {
		log.Fatalf("Failed to migrate %v", err)
	}
	// declare and initialise the app
	app := &App{
		router: gin.Default(),
		db:     db,
		routes: routes.NewRoutes(),
	}
	app.setupRoutes()
	return app, nil

}

func (a *App) setupRoutes() {
	a.router.GET("/api/v1/health-check", a.healthCheck)
	a.routes.SetupRoutes(a.router)

}
func (a *App) Run() error {
	return a.router.Run(":8080")
}

func (a *App) healthCheck(c *gin.Context) {
	response := HealthCheckResponse{
		Status:  "OK",
		Message: "All Systems Narmal 🚀",
	}
	sqlDB, err := a.db.DB()

	if err != nil {

		response.Status = "Error"
		response.Message = "Database Connection Error"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// ping the database
	err = sqlDB.Ping()
	if err != nil {
		response.Status = "Error"
		response.Message = "Database Ping Failed"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusOK, response)
}
