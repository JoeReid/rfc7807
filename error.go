package rfc7807

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Cause       error
	Status      int
	Type        string
	Detail      string
	ExtraFields map[string]interface{}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s %s: %s: %s", http.StatusText(e.Status), e.Type, e.Detail, e.Cause.Error())
}

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
