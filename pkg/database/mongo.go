package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect to mongo database
// Accept *context.Context and string as parameter
// Return *mongo.Database
func ConnectMongo(ctx context.Context, mongoURL string, mongoDB string) *mongo.Database {
	// Connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal("Mongo database connection error ", err)
	}

	// Set database
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo database connection error ", err)
	}
	fmt.Println("Mongo database connection successfully")
	db := client.Database(mongoDB)
	return db
}
