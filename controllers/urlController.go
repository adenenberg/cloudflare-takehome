package controllers

import (
	// "cloudflare-takehome/db"
	"net/http"
)

// var client = db.Connect()

var PingEndpoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
})

var CreateURLEndpoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// var
})
