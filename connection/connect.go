package connection

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/joshua-takyi/todo/helpers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"  // Import the options package for MongoDB client options
	"go.mongodb.org/mongo-driver/mongo/readpref" // For ping verification
)

var (
	// Client is the MongoDB client used throughout the application
	Client *mongo.Client
)

// Init initializes a connection to MongoDB using environment variables
// Returns error if connection fails or environment is misconfigured
func Init() error {
	// Load environment variables from .env.local file
	if err := godotenv.Load(".env.local"); err != nil {
		// We log but don't return - missing .env file is ok if env vars are set another way
		fmt.Printf("Warning: Failed to load .env.local file: %v\n", err)
	}

	// Get MongoDB URI from environment variables
	mongodb_uri := os.Getenv("MONGODB_URI")
	if mongodb_uri == "" {
		return helpers.Error{
			Message: "Missing environment variable: MONGODB_URI. Please set it in your .env file.",
			Status:  500, // Internal Server Error is more appropriate for configuration issues
		}
	}

	// Create a context with timeout for database operations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure context resources are freed

	// Configure client options
	clientOptions := options.Client().ApplyURI(mongodb_uri)

	// Connect to MongoDB
	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return helpers.Error{
			Message: fmt.Sprintf("Error connecting to MongoDB: %v", err),
			Status:  500,
		}
	}

	// Ping the database to verify connection
	err = Client.Ping(ctx, readpref.Primary())
	if err != nil {
		// Attempt to clean up failed connection
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer disconnectCancel()

		// Ignore disconnect error - we're already handling the ping error
		_ = Client.Disconnect(disconnectCtx)

		return helpers.Error{
			Message: fmt.Sprintf("Failed to verify MongoDB connection with ping: %v", err),
			Status:  500,
		}
	}

	fmt.Println("Successfully connected to MongoDB")
	return nil // No error, connection successful
}

// Close properly disconnects from MongoDB
// Should be called when the application is shutting down
func Close() error {
	if Client == nil {
		// Nothing to close
		return nil
	}

	// Create a timeout context for disconnect operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Disconnect from MongoDB
	err := Client.Disconnect(ctx)
	if err != nil {
		return helpers.Error{
			Message: fmt.Sprintf("Error disconnecting from MongoDB: %v", err),
			Status:  500,
		}
	}

	fmt.Println("Successfully disconnected from MongoDB")
	return nil
}
