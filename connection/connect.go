package connection

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/joshua-takyi/todo/helpers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options" // Import the options package for MongoDB client options
)

var (
	Client *mongo.Client
)

func Init() error {
	if err := godotenv.Load(".env.local"); err != nil {
		return helpers.Error{
			Message: fmt.Sprintf("Error loading .env file: %v", err),
			Status:  500, // Internal Server Error is more appropriate for configuration issues
		}
	}
	mongodb_uri := os.Getenv("MONGODB_URI")
	if mongodb_uri == "" {
		return helpers.Error{
			Message: "Missing environment variable: MONGODB_URI. Please set it in your .env file.",
			Status:  500, // Internal Server Error is more appropriate for configuration issues
		}
	}

	// create a context Timeout of 10 seconds

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Use the options package to create client options
	clientOptions := options.Client().ApplyURI(mongodb_uri)
	var err error

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		cancel()
		return helpers.Error{
			Message: fmt.Sprintf("Error connecting to MongoDB: %v", err),
			Status:  500, // Internal Server Error is more appropriate for connection issues
		}
	}

	// ping the database to check if the connection is successful
	err = client.Ping(ctx, nil)
	if err != nil {
		cancel()
		return helpers.Error{
			Message: fmt.Sprintf("Error pinging MongoDB: %v", err),
			Status:  500, // Internal Server Error is more appropriate for connection issues
		}
	}

	// Log the success message and return nil for error
	fmt.Println("MongoDB connection established successfully")
	cancel() // Ensure the context is canceled to release resources
	return nil
}

func Close(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		return helpers.Error{
			Message: fmt.Sprintf("Error disconnecting from MongoDB: %v", err),
			Status:  500, // Internal Server Error is more appropriate for disconnection issues
		}
	}

	fmt.Println("MongoDB connection closed successfully")
	return nil
}
