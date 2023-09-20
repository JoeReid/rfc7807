package rfc7807_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/JoeReid/rfc7807"
)

func ExampleError() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rfcErr := &rfc7807.Error{
			Cause:  errors.New("JWT token exp claim is in the past"),
			Status: http.StatusUnauthorized,
			Type:   "https://example.net/errors/invalid-jwt",
			Detail: "The JWT token provided is not valid",
		}

		rfc7807.ErrorResponse(rfcErr, w, r)
	})

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))

	// Output: {"detail":"The JWT token provided is not valid","status":401,"type":"https://example.net/errors/invalid-jwt"}
	fmt.Println(w.Body.String())
}

func ExampleError_extraFields() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rfcErr := &rfc7807.Error{
			Cause:  errors.New("JWT token exp claim is in the past"),
			Status: http.StatusUnauthorized,
			Type:   "https://example.net/errors/invalid-jwt",
			Detail: "The JWT token provided is not valid",
			ExtraFields: map[string]interface{}{
				"traceId": "1234567890",
			},
		}

		rfc7807.ErrorResponse(rfcErr, w, r)
	})

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))

	// Output: {"detail":"The JWT token provided is not valid","status":401,"traceId":"1234567890","type":"https://example.net/errors/invalid-jwt"}
	fmt.Println(w.Body.String())
}
