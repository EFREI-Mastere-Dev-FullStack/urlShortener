package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"urlShortener/database"
)

var ctx = context.TODO()

func main() {
	dbConnection, err := database.NewMongoDBConnection("mongodb://localhost:27017", "urlshortener")
	if err != nil {
		panic("Impossible d'initialiser la base de données MongoDB: " + err.Error())
	}

	if err := dbConnection.Client.Ping(ctx, nil); err != nil {
		panic("Impossible de se connecter à la base de données MongoDB: " + err.Error())
	}

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	router.Run(":8080")
}
