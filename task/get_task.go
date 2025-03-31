package task

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joshua-takyi/todo/connection"
	"github.com/joshua-takyi/todo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// get all tasks
func GetTask(ctx *gin.Context) {
	start := time.Now()
	var tasks []primitive.M

	collection := connection.Client.Database("Go").Collection("tasks")
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to retrieve tasks: " + err.Error()})
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var task primitive.M
		if err := cursor.Decode(&task); err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to decode task: " + err.Error()})
			return
		}
		tasks = append(tasks, task)
	}

	ended := time.Since(start).Seconds()
	ctx.JSON(200, gin.H{
		"message":  "Tasks retrieved successfully",
		"duration": ended,
		"tasks":    tasks,
	})
}

// get by id
func GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id format"})
		return
	}

	filter := bson.M{"_id": parsedId}
	collection := connection.Client.Database("Go").Collection("tasks")

	// Attempt to find a single task in the "tasks" collection that matches the provided filter.
	// The result is decoded into the `task` variable of type `model.Task`.
	var task model.Task
	err = collection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to retrieve task: " + err.Error()})
	}

	ctx.JSON(200, gin.H{
		"message": "Task retrieved successfully",
		"task":    task,
	})
}
