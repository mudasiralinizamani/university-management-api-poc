package data

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbInstance() *mongo.Client {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error occurred while loading .env file")
	}

	url := os.Getenv("MONGODB_CONNECTION")

	client, err := mongo.NewClient(options.Client().ApplyURI(url))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully Connected to Mongodb Database")

	return client
}

var mongoClient *mongo.Client = DbInstance()

func openCollection(dbClient *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = dbClient.Database("UniversityManagementDB").Collection(collectionName)
	return collection
}

// Collection Variables
// These variables are used in throughout in the application - Mudasir Ali
var UserCollection *mongo.Collection = openCollection(mongoClient, "Users")
var FacultyCollection *mongo.Collection = openCollection(mongoClient, "Faculties")
var DepartmentCollection *mongo.Collection = openCollection(mongoClient, "Departments")
