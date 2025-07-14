package api

import (
	"MedSearch/app/database"
	"MedSearch/app/models"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Search(c *gin.Context) {
	searchStr := c.Query("text")
	searchStr = refineFullTextQuery(searchStr)

	filter := bson.M{
		"$text": bson.M{
			"$search": searchStr, // q là từ khóa nhập từ người dùng
		},
	}

	collection := database.DB.Collection("drugs")
	cursor, err := collection.Find(context.Background(), filter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(context.Background())

	var drugs []models.Drug
	for cursor.Next(context.Background()) {
		var drug models.Drug
		cursor.Decode(&drug)
		drugs = append(drugs, drug)
	}

	c.JSON(http.StatusOK, drugs)
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
