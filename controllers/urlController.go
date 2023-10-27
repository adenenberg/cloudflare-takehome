package controllers

import (
	"cloudflare-takehome/db"
	"cloudflare-takehome/handlers"
	"cloudflare-takehome/models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fatih/color"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var client = db.Connect()

var PingEndpoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
})

var CreateURLEndpoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var shortenedURL models.ShortenedURL
	err := json.NewDecoder(r.Body).Decode((&shortenedURL))
	if err != nil {
		color.Red("Failed to decode json: %s", err)
		//todo err handling
		return
	}
	//todo validate

	id, _ := gonanoid.New()
	shortenedURL.ID = id
	shortenedURL.CreationDate = primitive.NewDateTimeFromTime(time.Now().UTC())

	collection := client.Database("cloudflare").Collection("shortened_url")
	result, err := collection.InsertOne(context.Background(), shortenedURL)
	if err != nil {
		color.Red("Failed to insert into DB: %s", err)
		//todo err handling
		return
	}

	color.White("Inserted new short url: %s", result.InsertedID)
	handlers.SuccessResponse(shortenedURL.GenerateShortUrl(), w)
})
