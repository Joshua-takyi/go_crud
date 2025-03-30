package helpers

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
