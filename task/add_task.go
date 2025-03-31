package task

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joshua-takyi/todo/connection"
	"github.com/joshua-takyi/todo/helpers"
	"github.com/joshua-takyi/todo/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTask(ctx *gin.Context) {
	// Parse the incoming JSON request into a Task struct
	var task model.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate the task fields using the helper function
	if err := helpers.ValidateTask(task); err != nil {
		ctx.JSON(err.GetStatus(), gin.H{"error": err.Error()})
		return
	}

	// Set default values for the task
	task.ID = primitive.NewObjectID()
	task.Metadata.CreatedAt = time.Now()
	task.Metadata.UpdatedAt = time.Now()
	task.Completed = false

	// Create a timeout context for database operations
	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure resources are released when function completes

	// Get a reference to the tasks collection
	collection := connection.Client.Database("Go").Collection("tasks")

	// Insert the new task into the database
	result, err := collection.InsertOne(dbCtx, task)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create task: " + err.Error()})
		return
	}

	// Return a success response with the created task
	ctx.JSON(201, gin.H{
		"message": "Task created successfully",
		"task":    task,
		"id":      result.InsertedID,
	})
}
