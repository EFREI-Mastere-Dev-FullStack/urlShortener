package router

import (
	"urlShortener/auth"
	"urlShortener/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("view/*")

	// Certains navigateur peuvent garder en cache l'url et ne plus incrementer le count
	// https://developer.mozilla.org/fr/docs/Web/HTTP/Headers/Cache-Control
	router.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache")
	})

	// Liste des routes :
	// Definir les methods dans controller/controller.go
	router.GET("/home", auth.RequiredAuth, controller.IndexPage)
	router.POST("/shorten", controller.ShortenURL)
	router.GET("/:slug", controller.RedirectURL)
	router.GET("/", controller.LoginPage)
	router.POST("/", controller.Login)
	router.GET("/register", controller.RegisterPage)
	router.POST("/register", controller.Register)
	router.GET("/logout", controller.Logout)
	return router
}
