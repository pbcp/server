package api

import (
	"github.com/alfredxing/clipboard/api/db"

	"net/http"
)

func Serve() {
	db.Setup()

	db := db.Open()
	defer db.Close()

	api := API{
		DB: db,
	}
	api.Handler = api.Routes()

	http.ListenAndServe(":9000", &api)
}
