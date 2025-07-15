package api

import (
	"MedSearch/app/database/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	searchTerm := c.Query("text")
	drugs, err := repository.SearchDrug(searchTerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}
	c.JSON(http.StatusOK, drugs)
}
