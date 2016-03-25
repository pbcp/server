package api

import (
	"github.com/julienschmidt/httprouter"
)

// Routes sets the routes and handlers for the API
func (a *API) Routes() *httprouter.Router {
	router := httprouter.New()

	// 'Cause I just wanna copy and paste
	router.GET("/:id", a.Paste)
	router.POST("/:id", a.Copy)

	// Get from history
	router.GET("/:id/:index", a.Retrieve)

	return router
}
