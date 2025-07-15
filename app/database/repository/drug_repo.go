package repository

import (
	"MedSearch/app/database"
	"MedSearch/app/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

func SearchDrug(term string) ([]models.Drug, error) {
	searchStr := refineFullTextQuery(term)

	query := elastic.NewMatchQuery("name", searchStr).Fuzziness("AUTO")

	searchResult, err := database.ElasticClient.Search().
		Index("drugs").
		Query(query).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	// Duyệt qua từng kết quả
	var drugs []models.Drug
	for _, hit := range searchResult.Hits.Hits {
		var d models.Drug
		err := json.Unmarshal(hit.Source, &d)
		if err != nil {
			fmt.Printf("Lỗi parse: %v\n", err)
			continue
		}

		// Lấy ID từ Elasticsearch (nếu không nằm trong source)

		drugs = append(drugs, d)
	}

	return drugs, nil
}

func refineFullTextQuery(query string) string {
	words := strings.Fields(query)
	var newQuery strings.Builder
	for _, word := range words {
		newQuery.WriteString("+")
		newQuery.WriteString(word)
	}

	return newQuery.String()
}
