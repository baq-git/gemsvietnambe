package errors

// Used for custom error
type Err struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (err Err) Error() string {
	return err.Message
}

// Sentinel errors
const ErrResourceUnavailable = "This resource is unavailable"
