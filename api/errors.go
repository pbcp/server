package api

import (
	"encoding/json"
	"net/http"
)

// Error is a of http.Handler used to serve error messages
type Error struct {
	Code    int
	Message string
}

func (e Error) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := map[string]string{
		"error": e.Message,
	}

	json, err := json.Marshal(res)
	if err != nil {
		json = []byte("{\"error\":\"Internal error\"}")
	}

	http.Error(w, string(json), e.Code)
}

// NotFound is an Error for service 404 errors.
var NotFound = Error{
	Code:    http.StatusNotFound,
	Message: "Not found",
}

// Unauthorized is an Error served when authentication fails.
var Unauthorized = Error{
	Code:    http.StatusUnauthorized,
	Message: "Authentication failed",
}

// BadRequest is an Error given to malformed requests.
var BadRequest = Error{
	Code:    http.StatusBadRequest,
	Message: "Malformed request or request parameters",
}

// InternalError is returned when the server encounters an error not likely caused by
// faults in the client's request.
var InternalError = Error{
	Code:    http.StatusInternalServerError,
	Message: "Internal server error encountered",
}
