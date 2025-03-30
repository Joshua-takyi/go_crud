package router

import "github.com/gin-gonic/gin"

// Package router provides the HTTP routes for the task management API.
// It defines the routes for creating, retrieving, updating, and deleting tasks.

func Router() *gin.Engine {
	router := gin.Default()

	// Define the routes for the task management API
	// v1 := router.Group("/api/v1")
	{
		// v1.POST("/tasks", createTask)       // Create a new task
		// v1.GET("/tasks", getTasks)          // Retrieve all tasks
		// v1.GET("/tasks/:id", getTask)       // Retrieve a specific task by ID
		// v1.PATCH("/tasks/:id", updateTask)  // Update a specific task by ID
		// v1.DELETE("/tasks/:id", deleteTask) // Delete a specific task by ID
	}

	return router
}
