package database

import (
	"MedSearch/app/logger"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var ElasticClient *elasticsearch.Client
var Logger *logger.MLogger

//var ElasticClient

func SetLogger(logger *logger.MLogger) {
	Logger = logger
}

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

	log.Println("✅ Connected to MongoDB!")

	ElasticClient, err = ConnectElasticsearch()
	if err != nil {
		Logger.Log("%s", err.Error())
	}

	time.Sleep(5 * time.Second)

	MigrateToElastic()
}

func ConnectElasticsearch() (*elasticsearch.Client, error) {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://localhost:9200"
	}

	var client *elasticsearch.Client
	var err error
	maxAttempts := 10
	for i := 1; i <= maxAttempts; i++ {
		// --- Kết nối tới Elasticsearch ---
		client, err = elasticsearch.NewDefaultClient()
		if err != nil {
			Logger.Log("Elasticsearch not ready (attempt %d/%d): %v", i, maxAttempts, err)
			time.Sleep(5 * time.Second)

		} else {
			break
		}

	}
	if err != nil {
		return nil, fmt.Errorf("error creating Elasticsearch client after retries: %w", err)
	}
	Logger.Log("Elasticsearch connected!")
	return client, nil
}
