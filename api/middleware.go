package api

import (
	"net/http"
)

// ServerHeader is the Server header value.
const ServerHeader string = "pbcp"

// ContentType is the Content-Type header for data responses.
const ContentType string = "application/json"

// SetHeaders sets default headers on all returned requests.
func (_ *API) SetHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", ServerHeader)
	w.Header().Set("Content-Type", ContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

// ServeHTTP serves as middleware for all requests to the API. It handles setting
// custom headers for all requests, and authenticating users.
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.SetHeaders(w, r)

	// Make sure endpoint exists
	handler, _, _ := a.Handler.Lookup(r.Method, r.URL.String())
	if r.Method != "OPTIONS" && handler == nil {
		NotFound.ServeHTTP(w, r)
		return
	}

	a.Handler.ServeHTTP(w, r)
}
