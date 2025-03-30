package response

import (
	"encoding/json"
)

type Error interface {
	Error() string
}

var _ Error = new(ErrorResponse)
var _ error = new(ErrorResponse)

// ErrorDetail untuk error spesifik dalam request
type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

func (e *ErrorResponse) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// SendError mengirimkan response error dengan format standar
func SendError(code, message string, details ...ErrorDetail) error {
	return &ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	}
}
