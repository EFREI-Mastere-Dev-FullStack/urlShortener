package router

import (
	"github.com/gin-gonic/gin"
	"urlShortener/controller"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Liste des routes :
	// Definir les methods dans controller/controller.go
	router.GET("/", controller.IndexPage)

	return router
}
