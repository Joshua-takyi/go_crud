package task

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joshua-takyi/todo/connection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PatchTask handles partial updates to a task document identified by ID
// It validates the input, processes the update, and returns appropriate responses
func PatchTask(ctx *gin.Context) {
	// Extract and validate the task ID from URL parameters
	paramsId := ctx.Param("id")
	id, err := primitive.ObjectIDFromHex(paramsId)
	if err != nil {
		// Return a structured error response with HTTP 400 Bad Request
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid task ID format",
			"details": fmt.Sprintf("The provided ID '%s' is not a valid MongoDB ObjectID", paramsId),
		})
		return // Added return statement to prevent execution continuing after error
	}

	// Parse the request body into a map for flexible field updates
	var updateFields map[string]interface{}
	if err := ctx.ShouldBindJSON(&updateFields); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Validate update payload - cannot be empty
	if len(updateFields) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Update payload cannot be empty",
		})
		return
	}

	// Validate that protected fields are not being modified
	protectedFields := []string{"id", "_id", "created_at"}
	for _, field := range protectedFields {
		if _, exists := updateFields[field]; exists {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Protected field modification attempt",
				"details": fmt.Sprintf("The field '%s' cannot be updated", field),
			})
			return
		}
	}

	// Type-specific validation for common fields
	if completed, exists := updateFields["completed"]; exists {
		// Ensure completed is a boolean if present
		if _, ok := completed.(bool); !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid field type",
				"details": "Field 'completed' must be a boolean value",
			})
			return
		}
	}

	// Add update timestamp
	updateFields["updated_at"] = time.Now()

	// Prepare the MongoDB filter and update operation
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": updateFields,
	}

	// Create a context with timeout for database operations
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Ensure resources are freed

	// Get database collection
	collection := connection.Client.Database("Go").Collection("tasks")

	// Execute the update operation
	result, err := collection.UpdateOne(dbCtx, filter, update)
	if err != nil {
		// Handle database operation errors
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database operation failed",
			"details": err.Error(),
		})
		return
	}

	// Check if the task was found
	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Task not found",
			"details": fmt.Sprintf("No task exists with ID: %s", paramsId),
		})
		return
	}

	// Fetch the updated task to return complete data
	var updatedTask bson.M
	err = collection.FindOne(dbCtx, filter).Decode(&updatedTask)
	if err != nil {
		// Log the error but still return success since the update worked
		fmt.Printf("Warning: Update succeeded but failed to fetch updated task: %v\n", err)
		// Return partial success response
		ctx.JSON(http.StatusOK, gin.H{
			"message":  "Task updated successfully",
			"modified": updateFields,
			"warning":  "Could not retrieve the complete updated task",
		})
		return
	}

	// Return success response with the complete updated task
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"task":    updatedTask,
	})
}
