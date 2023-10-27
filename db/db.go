package db

import (
	"context"
	"log"

	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var client *mongo.Client

func Connect() *mongo.Client {
	ctx := context.Background()

	clientOptions := options.Client().ApplyURI("mongodb://root:cloudflare@localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Connection Failed to Database")
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Connection Failed to Database")
		log.Fatal(err)
	}

	// defer func() {
	// 	if err := client.Disconnect(ctx); err != nil {
	// 		log.Fatalf("mongodb disconnect error : %v", err)
	// 	}
	// }()

	color.Green("Connected to Database")
	return client
}
