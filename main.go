package main

import (
	"context"
	"fmt"
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

	fmt.Println("Connecté à MongoDB")
}
