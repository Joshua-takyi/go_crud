package helpers

import "github.com/joshua-takyi/todo/model"

// Response represents a successful operation with a message and status code
type Response struct {
	Message string
	Status  int
}

// Error represents an error condition with a message and status code
// It implements the error interface
type Error struct {
	Message string
	Status  int
}

// Success returns the success message
func (r Response) Success() string {
	return r.Message
}

// Error implements the error interface to allow using Error as a standard Go error
func (e Error) Error() string {
	return e.Message
}

// GetStatus returns the HTTP status code associated with this error
// This allows callers to extract the status code separately from the error message
func (e Error) GetStatus() int {
	return e.Status
}

func ValidateTask(task model.Task) *Error {
	// Validate each field of the task and return an error if a field is missing
	if task.Title == "" {
		return &Error{Message: "Title is required", Status: 400}
	}
	if task.Description == "" {
		return &Error{Message: "Description is required", Status: 400}
	}
	if task.Priority == "" {
		return &Error{Message: "Priority is required", Status: 400}
	}
	if len(task.Tags) == 0 {
		return &Error{Message: "Tags are required", Status: 400}
	}
	if len(task.Image) == 0 {
		return &Error{Message: "Image is required", Status: 400}
	}

	return nil
}
