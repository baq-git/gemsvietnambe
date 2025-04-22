package errors

// Used for custom error
type Err struct {
	Message  string            `json:"message,omitempty"`
	Messages map[string]string `json:"messages,omitempty"`
	Data     any               `json:"data,omitempty"`
}

func (err Err) Error() string {
	return err.Message
}

func (err Err) Errors() map[string]string {
	return err.Messages
}

// Sentinel errors
const ErrResourceUnavailable = "This resource is unavailable"

// Validates Errors
const ErrValidateFailed = "Validate failed, check client parameters"
