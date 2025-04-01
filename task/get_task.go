package task

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joshua-takyi/todo/connection"
	"github.com/joshua-takyi/todo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetTask retrieves tasks with pagination support
// Query parameters:
// - page: current page number (default: 1)
// - limit: number of tasks per page (default: 10)
func GetTask(ctx *gin.Context) {
	start := time.Now()

	// Parse pagination parameters
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10 // Default limit with a reasonable maximum
	}

	// Calculate skip value for pagination
	skip := (page - 1) * limit

	// Prepare options for MongoDB query
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.M{"createdAt": -1}) // Sort by creation date, newest first

	// Initialize an empty slice to store the tasks
	var tasks []primitive.M
	collection := connection.Client.Database("Go").Collection("tasks")

	// Count total documents for pagination info
	total, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to count tasks: " + err.Error()})
		return
	}

	// Execute the find query with options
	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to retrieve tasks: " + err.Error()})
		return // Important: return after error response
	}
	defer cursor.Close(context.Background())

	// Iterate through the cursor to decode documents
	for cursor.Next(context.Background()) {
		var task primitive.M
		if err := cursor.Decode(&task); err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to decode task: " + err.Error()})
			return
		}
		tasks = append(tasks, task)
	}

	// Check for cursor iteration errors
	if err := cursor.Err(); err != nil {
		ctx.JSON(500, gin.H{"error": "Cursor error: " + err.Error()})
		return
	}

	// Calculate pagination metadata
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Calculate execution time
	ended := time.Since(start).Seconds()

	// Return the paginated results with metadata
	ctx.JSON(200, gin.H{
		"message":  "Tasks retrieved successfully",
		"duration": ended,
		"tasks":    tasks,
		"pagination": gin.H{
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": totalPages,
			"hasMore":    page < totalPages,
		},
	})
}

// GetById retrieves a specific task by its ID
func GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	filter := bson.M{"_id": parsedId}
	collection := connection.Client.Database("Go").Collection("tasks")

	// Attempt to find a single task in the "tasks" collection that matches the provided filter.
	var task model.Task
	err = collection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Task not found: " + err.Error()})
		return // Important: return after error response
	}

	ctx.JSON(200, gin.H{
		"message": "Task retrieved successfully",
		"task":    task,
	})
}
