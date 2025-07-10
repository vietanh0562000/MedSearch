package routes

import (
	"MedSearch/app/api"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {
	apiRouter := router.Group("/v1/api")
	{
		apiRouter.GET("/search", api.Search)
	}
}
