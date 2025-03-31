ackage connection

import (
	"context"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// TestInit tests the Init function for MongoDB connection
func TestInit(t *testing.T) {
	// Load environment variables for testing
	if err := godotenv.Load(".env.local"); err != nil {
		t.Fatalf("Failed to load .env.local file: %v", err)
	}

	// Initialize the connection
	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Verify the connection with a ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx, readpref.Primary()); err != nil {
		t.Errorf("Ping to MongoDB failed: %v", err)
	}

	// Close the connection after test
	if err := Close(); err != nil {
		t.Errorf("Close() failed: %v", err)
	}
}