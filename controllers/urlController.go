package controllers

import (
	"cloudflare-takehome/db"
	"cloudflare-takehome/handlers"
	"cloudflare-takehome/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client = db.Connect()

const dbName = "cloudflare"
const urlTable = "shortened_url"
const statsTable = "url_stats"

var PingEndpoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
})

var CreateURLEndpoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var shortenedURL models.ShortenedURL
	err := json.NewDecoder(r.Body).Decode((&shortenedURL))
	if err != nil {
		color.Red("Failed to decode json: %s", err)
		handlers.ErrorResponse("Request data formatted incorrectly: "+err.Error(), w)
	}

	if shortenedURL.OriginalURL == "" {
		color.Red("Missing URL to shorten: %s", err)
		handlers.ErrorResponse("Must include URL", w)
	}

	u, _ := url.Parse(shortenedURL.OriginalURL)
	if u.Scheme == "" {
		shortenedURL.OriginalURL = "http://" + shortenedURL.OriginalURL
	}

	id, _ := gonanoid.New()
	shortenedURL.ID = id
	shortenedURL.CreationDate = primitive.NewDateTimeFromTime(time.Now().UTC())

	collection := client.Database(dbName).Collection(urlTable)
	result, err := collection.InsertOne(context.Background(), shortenedURL)
	if err != nil {
		color.Red("Failed to insert into DB: %s", err)
		handlers.ServerErrorResponse("Failed to insert into DB: "+err.Error(), w)
	}

	color.White("Inserted new short url: %s", result.InsertedID)
	resp := map[string]string{
		"short_url": shortenedURL.GenerateShortUrl(),
		"key":       shortenedURL.ID,
	}
	handlers.SuccessResponse(resp, w)
})

var GoToURLEndpoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var shortenedURL models.ShortenedURL

	idFilter := bson.D{primitive.E{Key: "_id", Value: params["id"]}}

	collection := client.Database(dbName).Collection(urlTable)
	err := collection.FindOne(context.Background(), idFilter).Decode(&shortenedURL)
	if err != nil {
		color.Red("Record not found: %s", err)
		handlers.ErrorResponse("Record not found", w)
	}

	//Check if URL has expired
	now := time.Now().UTC()
	if shortenedURL.ExpirationDate != 0 && now.After(shortenedURL.ExpirationDate.Time()) {
		color.Red("Record not found: %s", err)
		//Delete if expired
		collection.DeleteOne(context.Background(), idFilter)
		handlers.ErrorResponse("Record not found", w)
	}

	statsCollection := client.Database(dbName).Collection(statsTable)
	statsCollection.UpdateByID(context.Background(), params["id"],
		bson.M{"$push": bson.M{
			"access_times": primitive.NewDateTimeFromTime(now),
		}},
		options.Update().SetUpsert(true))

	color.White("Redirecting to URL")
	http.Redirect(w, r, shortenedURL.OriginalURL, http.StatusTemporaryRedirect)
})

var DeleteURLEndpoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var shortenedURL models.ShortenedURL

	idFilter := bson.D{primitive.E{Key: "_id", Value: params["id"]}}

	collection := client.Database("cloudflare").Collection("shortened_url")
	err := collection.FindOne(context.Background(), idFilter).Decode(&shortenedURL)
	if err != nil {
		color.Red("Record not found: %s", err)
		handlers.ErrorResponse("Record not found", w)
	}

	_, err = collection.DeleteOne(context.Background(), idFilter)
	if err != nil {
		color.Red("Record could not be deleted: %s", err)
		handlers.ServerErrorResponse("Error deleting record: "+err.Error(), w)
	}

	statsCollection := client.Database(dbName).Collection(statsTable)
	statsCollection.DeleteOne(context.Background(), idFilter)

	color.White("URL deleted: " + shortenedURL.ID)
	handlers.SuccessResponse("deleted", w)
})

var URLStatsEndpoint = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var urlStats models.URLStats

	idFilter := bson.D{primitive.E{Key: "_id", Value: params["id"]}}

	statsCollection := client.Database(dbName).Collection(statsTable)
	err := statsCollection.FindOne(context.Background(), idFilter).Decode(&urlStats)
	if err != nil {
		color.Red("Record not found: %s", err)
		handlers.ErrorResponse("Record not found", w)
	}

	now := time.Now().UTC()
	day := now.AddDate(0, 0, -1)
	week := now.AddDate(0, 0, -7)
	dayCount := 0
	weekCount := 0
	allCount := 0

	//todo: try this with a query instead
	for _, d := range urlStats.AccessTimes {
		if d.Time().After(day) {
			dayCount++
		} else if d.Time().After(week) {
			weekCount++
		}
		allCount++
	}

	result := map[string]int{
		"past_day":  dayCount,
		"past_week": weekCount + dayCount,
		"all_time":  allCount,
	}
	fmt.Println(result)

	handlers.SuccessResponse(result, w)
})
