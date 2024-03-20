package router

import (
	"github.com/gin-gonic/gin"
	"urlShortener/controller"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("view/*")

	// Liste des routes :
	// Definir les methods dans controller/controller.go
	router.GET("/", controller.IndexPage)
	router.POST("/shorten", controller.ShortenURL)

	return router
}
