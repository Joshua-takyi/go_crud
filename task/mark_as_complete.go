package task

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshua-takyi/todo/connection"
	"github.com/joshua-takyi/todo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MarkAsComplete(ctx *gin.Context) {
	// Get task ID from URL parameter
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	// Convert string ID to MongoDB ObjectID
	parsedId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Initialize database collection
	collection := connection.Client.Database("Go").Collection("tasks")

	// First, find the current task to check its completion status
	var task model.Task
	filter := bson.M{"_id": parsedId}
	err = collection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Task not found",
			"details": err.Error(),
		})
		return
	}

	// Toggle the completion status (true becomes false, false becomes true)
	newCompletionStatus := !task.Completed

	// Prepare update operation with the new status
	update := bson.M{"$set": bson.M{"completed": newCompletionStatus}}

	// Update the task with the new completion status
	result := collection.FindOneAndUpdate(context.Background(), filter, update)
	if err := result.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update task completion status",
			"details": err.Error(),
		})
		return
	}

	// Prepare response message based on the new status
	var message string
	if newCompletionStatus {
		message = "Task marked as complete"
	} else {
		message = "Task marked as incomplete"
	}

	// Return success response with appropriate message
	ctx.JSON(http.StatusOK, gin.H{
		"message":   message,
		"task_id":   id,
		"completed": newCompletionStatus,
	})
}
