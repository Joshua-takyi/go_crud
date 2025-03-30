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

	r := router.Router()
	fmt.Println("server starting on port 8080")
	if err := r.Run(":8080"); err != nil {
		// if the server fails to start we return an error message and exit the program
		fmt.Println(err.Error())
		return
	}
	defer func() {
		if err := connection.Close(connection.Client); err != nil {
			fmt.Println(err.Error())
		}
	}()

	// if the server starts successfully we return a success message
	fmt.Println("server started successfully on port 8080")

}
