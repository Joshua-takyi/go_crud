package main

import (
	"fmt"
	"os"

	"github.com/joshua-takyi/todo/connection"
	"github.com/joshua-takyi/todo/router"
)

func main() {
	if err := connection.Init(); err != nil {
		fmt.Println("Database connection error:", err.Error())
		return
	}

	defer func() {
		if err := connection.Close(); err != nil {
			fmt.Println("Error closing database connection:", err.Error())
		}
	}()

	r := router.Router()

	// Get port from environment variable for cloud deployment compatibility
	// Platforms like Render typically provide the port via the PORT environment variable
	port := os.Getenv("PORT")
	if port == "" {
		// Default to port 8080 for local development
		port = "8080"
	}

	// Log server startup information
	fmt.Printf("Server starting on port %s\n", port)

	// Start the HTTP server with the configured port
	// Run() blocks until the server is stopped or an error occurs
	if err := r.Run(":" + port); err != nil {
		fmt.Println("Server startup error:", err.Error())
		return
	}
}
