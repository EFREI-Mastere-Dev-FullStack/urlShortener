package main

import (
	"context"
	"urlShortener/database"
	"urlShortener/router"
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

	r := router.SetupRouter()

	r.LoadHTMLGlob("view/*")
	r.Run(":8080")
}
