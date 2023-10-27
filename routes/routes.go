package routes

import (
	"cloudflare-takehome/controllers"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", controllers.PingEndpoint).Methods("GET")
	router.HandleFunc("/create", controllers.CreateURLEndpoint).Methods("POST")
	router.HandleFunc("/go/{id}", controllers.GoToURLEndpoint).Methods("GET")
	router.HandleFunc("/delete/{id}", controllers.DeleteURLEndpoint).Methods("DELETE")
	router.HandleFunc("/stats/{id}", controllers.URLStatsEndpoint).Methods("GET")

	return router
}
