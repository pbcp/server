package api

import (
	"github.com/julienschmidt/httprouter"
	"os"
)

// Routes sets the routes and handlers for the API
func (a *API) Routes() *httprouter.Router {
	router := httprouter.New()

	// 'Cause I just wanna copy and paste
	router.POST("/board/:id", a.Copy)
	router.GET("/board/:id/:index", a.Paste)

	router.GET("/register", a.Register)

	if os.Getenv("DEGUG") != "" {
		router.GET("/debug", a.Debug)
	}

	return router
}
