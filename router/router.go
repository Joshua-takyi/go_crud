package router

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/joshua-takyi/todo/task"
)

// Package router provides the HTTP routes for the task management API.
// It defines the routes for creating, retrieving, updating, and deleting tasks.

func Router() *gin.Engine {
	router := gin.Default()

	// Load environment variables from .env.local file
	if err := godotenv.Load(".env.local"); err != nil {
		fmt.Printf("Warning: Failed to load .env.local file: %v\n", err)
	}

	// Get frontend URL from environment variables
	frontendUrl := os.Getenv("FRONTEND_URL")
	if frontendUrl == "" {
		frontendUrl = "http://localhost:3000" // Default to localhost if not set
		fmt.Println("FRONTEND_URL not set, using default:", frontendUrl)
	}

	// Configure CORS middleware with proper settings
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendUrl},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum cache time for preflight requests (in seconds)
	}))

	// Define the routes for the task management API
	v1 := router.Group("/api/v1")
	{
		v1.POST("/tasks", task.CreateTask)       // Create a new task
		v1.GET("/tasks", task.GetTask)           // Retrieve all tasks
		v1.GET("/tasks/:id", task.GetById)       // Retrieve a specific task by ID
		v1.PATCH("/tasks/:id", task.PatchTask)   // Update a specific task by ID
		v1.DELETE("/tasks/:id", task.DeleteTask) // Delete a specific task by ID
		// mark as complete
		v1.PATCH("/tasks/:id/complete", task.MarkAsComplete) // Mark a task as complete
	}

	return router
}
