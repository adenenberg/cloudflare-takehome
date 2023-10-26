package routes

import (
	"cloudflare-takehome/controllers"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", controllers.PingEndpoint).Methods("GET")

	return router
}
