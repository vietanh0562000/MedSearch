package repository

import (
	"MedSearch/app/database"
	"MedSearch/app/models"
	"bytes"
	"context"
	"encoding/json"
	"strings"

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

// Define a struct for the hits:
type Hit struct {
	Source models.Drug `json:"_source"`
}
type SearchHits struct {
	Hits struct {
		Hits []Hit `json:"hits"`
	} `json:"hits"`
}

func SearchDrug(term string) ([]models.Drug, error) {
	// Build the query as a JSON string
	query := `{
	"query": {
		"match": {
		"name": {
			"query": "` + term + `",
			"fuzziness": "AUTO"
		}
		}
	}
	}`

	// Perform the search
	res, err := database.ElasticClient.Search(
		database.ElasticClient.Search.WithContext(context.Background()),
		database.ElasticClient.Search.WithIndex("drugs"),
		database.ElasticClient.Search.WithBody(bytes.NewReader([]byte(query))),
		database.ElasticClient.Search.WithTrackTotalHits(true),
		database.ElasticClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// Then unmarshal:
	var sh SearchHits
	if err := json.NewDecoder(res.Body).Decode(&sh); err != nil {
		return nil, err
	}

	var drugs []models.Drug
	for _, hit := range sh.Hits.Hits {
		drugs = append(drugs, hit.Source)
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
