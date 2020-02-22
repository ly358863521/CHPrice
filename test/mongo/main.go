package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	provinceSET := client.Database("CHPrice").Collection("province")
	var province bson.M
	// filter := bson.D{primitive.E{Key: "name", Value: primitive.Regex{Pattern: "北京", Options: ""}}}
	// err = provinceSET.FindOne(context.Background(), filter).Decode(&province)
	err = provinceSET.FindOne(context.Background(), bson.D{{Key: "name", Value: "四川省"}, {Key: "price", Value: 6702}}).Decode(&province)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(province)
}
