package app

import (
	"MedSearch/app/config"
	"MedSearch/app/database"
	"MedSearch/app/database/repository"
	"MedSearch/app/logger"
	"MedSearch/app/routes"
	"fmt"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func Start(config *config.AppConfig, logger *logger.MLogger) {
	database.Connect(config.GetDbURI(), config.GetDbName())
	repository.CreateDrugTextIndex()

	router := gin.Default()

	// Add CORS middleware with simpler configuration
	router.Use(cors.Default())

	routes.Setup(router)

	address := fmt.Sprintf(":%s", config.GetPort())
	router.Run(address)
}
