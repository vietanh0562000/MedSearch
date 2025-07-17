package repository

import (
	"MedSearch/app/database"
	"MedSearch/app/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	if res.IsError() {
		log.Fatalf("Search failed: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("Error parsing response: %s", err)
	}

	var drugIds []string

	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		doc := hit.(map[string]interface{})
		drugIds = append(drugIds, doc["_id"].(string))
	}
	// Then unmarshal:
	// var sh SearchHits
	// if err := json.NewDecoder(res.Body).Decode(&sh); err != nil {
	// 	return nil, err
	// }

	var drugs []models.Drug
	collection := database.DB.Collection("drugs")

	foundData, e := FindByIDs(context.Background(), collection, drugIds)
	if e != nil {
		return nil, e
	}

	drugs, err = ConvertBsonMToStruct(foundData)
	if err != nil {
		return nil, err
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

func FindByIDs(ctx context.Context, collection *mongo.Collection, ids []string) ([]bson.M, error) {
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, objID)
	}

	fmt.Println(objectIDs)

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func ConvertBsonMToStruct(bsonList []bson.M) ([]models.Drug, error) {
	var drugs []models.Drug
	for _, item := range bsonList {
		// Marshal from bson to []byte
		jsonBytes, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}

		// parse from []byte to Drug
		var drug models.Drug
		if err := json.Unmarshal(jsonBytes, &drug); err != nil {
			return nil, err
		}
		drugs = append(drugs, drug)
	}
	return drugs, nil
}
