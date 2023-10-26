package main

import (
	"cloudflare-takehome/routes"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/rs/cors"
	// "github.com/matoous/go-nanoid"
)

func main() {
	color.Cyan("Server running on localhost: 8080")

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	router := routes.Routes()
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	handler := c.Handler(router)
	http.ListenAndServe(":8080", handler)
}

// This function ensures that a number of keys are always available for the service to use.
// When a certain threshold is met, it will generate a new set of keys to be used.
// The keys will be stored in a table to be used by the API
func generateKeys() {
	//TODO: use nanoid
}
