package rfc7807

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error is an implementation of RFC7807. It is a wrapper around an error that
// provides a status code, type, and detail. It also allows for extra fields to
// be added to the response.
//
// Error implements the error interface, so it can be used anywhere an error is
// expected. It also implements json.Marshaler, so it can be marshaled to JSON
// and written to the client, however, ErrorResponse should be used for this purpose.
type Error struct {
	Cause       error
	Status      int
	Type        string
	Detail      string
	ExtraFields map[string]interface{}
}

// Error returns a string representation of the error.
func (e *Error) Error() string {
	return fmt.Sprintf("%s %s: %s: %s", http.StatusText(e.Status), e.Type, e.Detail, e.Cause.Error())
}

// MarshalJSON marshals the error to JSON, popuilating zero values with sane defaults.
func (e *Error) MarshalJSON() ([]byte, error) {
	if e.Status == 0 {
		e.Status = http.StatusInternalServerError
	}

	if e.Type == "" {
		e.Type = "unknown-error"
	}

	if e.Detail == "" && e.Cause != nil {
		e.Detail = e.Cause.Error()
	}

	fields := map[string]interface{}{
		"status": e.Status,
		"type":   e.Type,
		"detail": e.Detail,
	}

	for k, v := range e.ExtraFields {
		fields[k] = v
	}

	return json.Marshal(fields)
}

// ErrorResponse is a helper function to write an RFC7807 error response to the
// client. If the error is not an RFC7807 error, it will be wrapped in one with
// a status of 500 and a type of "unknown-error".
//
// If the error cannot be marshaled to JSON, a 500 will be returned with a
// message indicating that the error could not be encoded.
//
// Once the function returns, the response will have been written to and the
// request and callers should also return.
func ErrorResponse(err error, w http.ResponseWriter, r *http.Request) {
	rfcErr, ok := err.(*Error)
	if !ok {
		rfcErr = &Error{
			Status: http.StatusInternalServerError,
			Type:   "unknown-error",
			Cause:  err,
		}
	}

	data, err := json.Marshal(rfcErr)
	if err != nil {
		errJSON := fmt.Sprintf(`{"status": %d, "type": "%s", "detail": "%s"}`,
			http.StatusInternalServerError,
			"rfc7807-encoding-error",
			err.Error(),
		)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errJSON))
		return
	}

	w.WriteHeader(rfcErr.Status)
	w.Write(data)
}
