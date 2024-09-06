package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	connect        *mongo.Client
	Collection     *mongo.Collection
	UserCollection *mongo.Collection
)

func DbConnect() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error occured while loading env")
	}
	URL := os.Getenv("MONGO_URL")
	DATABASE := os.Getenv("DATABASE")
	COLLECTION := os.Getenv("COLLECTION")
	USERCOLLECTION := os.Getenv("USERCOLLECTION")
	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(URL).SetServerAPIOptions(serverApi)
	client, err := mongo.Connect(context.Background(), opts)
	/* connect = client */
	Collection = client.Database(DATABASE).Collection(COLLECTION)
	UserCollection = client.Database(DATABASE).Collection(USERCOLLECTION) //*****add collection*****
	if err != nil {
		log.Fatal(err)
	}
	err = client.Database("big_tech").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err()
	return err
}
func Close() error {
	return connect.Disconnect(context.Background())
}
