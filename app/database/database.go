package database

import (
	"MedSearch/app/database/repository"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func Connect(uri string, dbName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		{
			log.Fatal(err)
		}
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	DB = client.Database(dbName)

	err = repository.CreateDrugTextIndex()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("✅ Connected to MongoDB!")
}
