package main

import (
	"cloudflare-takehome/routes"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/rs/cors"
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
