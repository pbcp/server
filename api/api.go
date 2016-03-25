package api

import (
	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
)

// API is the router struct to which all handlers are attached. It contains the global
// database connection, which is kept open and shared amongst all requests.
type API struct {
	Handler *httprouter.Router
	DB      *bolt.DB
}
