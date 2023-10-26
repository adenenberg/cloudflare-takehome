package controllers

import (
	// "cloudflare-takehome/db"
	"net/http"
)

// var client = db.Connect()

var PingEndpoint = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("pong"))
})
