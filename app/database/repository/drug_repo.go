package repository

import (
	"MedSearch/app/database"
	"MedSearch/app/models"
	"context"
)

func InsertDrug(drug *models.Drug) error {
	collection := database.DB.Collection("drugs")
	_, err := collection.InsertOne(context.Background(), drug)
	return err
}
