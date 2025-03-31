package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/joshua-takyi/todo/task"
)

func Router() *gin.Engine {
	router := gin.Default()

	if err := godotenv.Load("../.env.local"); err != nil {
		fmt.Printf("Warning: Failed to load .env.local file: %v\n", err)
		fmt.Println("This is expected in production environments where environment variables are set differently")
	}

	frontendUrl := os.Getenv("FRONTEND_URL")
	if frontendUrl == "" {
		frontendUrl = "http://localhost:3000" // Default to localhost if not set
		fmt.Println("FRONTEND_URL not set, using default:", frontendUrl)
	}

	// Configure CORS middleware with proper settings to allow cross-origin requests
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendUrl}, // Allow the frontend URL and all origins for testing
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum cache time for preflight requests (in seconds)
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Todo API is running",
			"version": "1.0",
			"endpoints": []string{
				"/api/v1/tasks - GET, POST",
				"/api/v1/tasks/:id - GET, PATCH, DELETE",
				"/api/v1/tasks/:id/complete - PATCH",
			},
		})
	})

	// Add a health check endpoint for monitoring systems
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Define the routes for the task management API under /api/v1 prefix
	v1 := router.Group("/api/v1")
	{
		v1.POST("/tasks", task.CreateTask)                   // Create a new task
		v1.GET("/tasks", task.GetTask)                       // Retrieve all tasks
		v1.GET("/tasks/:id", task.GetById)                   // Retrieve a specific task by ID
		v1.PATCH("/tasks/:id", task.PatchTask)               // Update a specific task by ID
		v1.DELETE("/tasks/:id", task.DeleteTask)             // Delete a specific task by ID
		v1.PATCH("/tasks/:id/complete", task.MarkAsComplete) // Mark a task as complete
	}

	return router
}
