package main

import (
	"log"

	"web-app-uji-emisi/config"
	"web-app-uji-emisi/controllers"
	"web-app-uji-emisi/routes"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to database
	config.ConnectDatabase()

	// Seed initial data
	controllers.SeedInitialData()

	// Setup routes
	r := routes.SetupRoutes()

	// Start server
	log.Printf("Server starting on port %s", config.AppConfig.Port)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
