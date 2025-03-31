package task

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joshua-takyi/todo/connection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteTask(ctx *gin.Context) {

	paramId := ctx.Param("id")

	if paramId == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Task not found",
			"error":   fmt.Sprintf("Task with ID %s not found", paramId),
		})
	}

	id, err := primitive.ObjectIDFromHex(paramId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid ID format",
			"error":   fmt.Sprintf("The provided ID '%s' is not a valid MongoDB ObjectID", paramId),
		})
		return
	}

	filter := bson.M{"_id": id}

	collection := connection.Client.Database("Go").Collection("tasks")

	_, err = collection.DeleteOne(context.Background(), filter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete task",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{
		"message": "Task deleted successfully",
	})
}
