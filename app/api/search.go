package api

import (
	"MedSearch/app/database"
	"MedSearch/app/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Search(c *gin.Context) {
	searchStr := c.Query("search")

	filter := bson.M{
		"$text": bson.M{
			"$search": searchStr, // q là từ khóa nhập từ người dùng
		},
	}

	opts := options.Find().SetLimit(int64(10)).SetSkip(int64((1 - 1) * 10))
	collection := database.DB.Collection("drugs")
	cursor, err := collection.Find(context.Background(), filter, opts)

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
