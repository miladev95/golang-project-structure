package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/yourproject/internal/config"
	"github.com/yourusername/yourproject/internal/di"
	"github.com/yourusername/yourproject/internal/di/modules"
	"github.com/yourusername/yourproject/internal/handlers/http/routes"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create DI container
	container := di.NewContainer()

	// Register modules
	container.
		RegisterModule(modules.NewUserModule())
	// Add more modules here as needed
	// .RegisterModule(modules.NewProductModule())
	// .RegisterModule(modules.NewOrderModule())

	// Setup dependencies
	if err := container.Setup(cfg); err != nil {
		log.Fatalf("Failed to setup dependencies: %v", err)
	}

	// Run database migrations
	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := config.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create Gin router
	router := gin.Default()

	// Get handlers from container
	userHandler, err := container.GetUserHandler()
	if err != nil {
		log.Fatalf("Failed to get user handler: %v", err)
	}

	// Register all routes
	routes.RegisterAll(
		router,
		routes.NewUserRouter(userHandler),
		// routes.NewProductRouter(productHandler), // Add more routers as needed
		// routes.NewOrderRouter(orderHandler),
	)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Start server
	log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := router.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}