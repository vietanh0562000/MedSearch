package repository

import (
	"MedSearch/app/database"
	"MedSearch/app/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertDrug(drug *models.Drug) error {
	collection := database.DB.Collection("drugs")
	_, err := collection.InsertOne(context.Background(), drug)
	return err
}

func CreateDrugTextIndex() error {
	collection := database.DB.Collection("drugs")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: "text"},
			{Key: "ingredients", Value: "text"},
			{Key: "uses", Value: "text"},
		},
		Options: options.Index().SetName("fulltext").SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}
