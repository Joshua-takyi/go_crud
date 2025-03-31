package router

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joshua-takyi/todo/task"
)

// Package router provides the HTTP routes for the task management API.
// It defines the routes for creating, retrieving, updating, and deleting tasks.

func Router() *gin.Engine {
	router := gin.Default()

	// cors
	frontendUrl := os.Getenv("FRONTEND_URL")
	if frontendUrl == "" {
		frontendUrl = "http://localhost:3000" // Default to localhost if not set
	}

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowOrigins:    []string{frontendUrl},
		AllowMethods:    []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:   []string{"Content-Length"},
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
