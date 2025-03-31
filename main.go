package main

import (
	"fmt"

	"github.com/joshua-takyi/todo/connection"
	"github.com/joshua-takyi/todo/router"
)

func main() {
	if err := connection.Init(); err != nil {
		// if connection fails we return an error message and exit the program
		fmt.Println(err.Error())
		return
	}

	// Ensure the MongoDB connection is closed when the application exits
	defer func() {
		if err := connection.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	r := router.Router()
	fmt.Println("server starting on port 8080")
	if err := r.Run(":8080"); err != nil {
		// if the server fails to start we return an error message and exit the program
		fmt.Println(err.Error())
		return
	}

	// Note: This line will never be reached because r.Run() blocks until the server is stopped
	fmt.Println("server started successfully on port 8080")
}
