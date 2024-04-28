package main

import (
	"log"
	"os"
	"urlShortener/database"
	"urlShortener/router"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)
	dbConnection, err := database.NewMongoDBConnection("mongodb://localhost:27017", "urlshortener")
	if err != nil {
		panic("Impossible d'initialiser la base de données MongoDB: " + err.Error())
	}

	database.Connection = dbConnection

	r := router.SetupRouter()

	err = r.Run(":8080")
	if err != nil {
		panic("Impossible démarrer le serveur: " + err.Error())
	}
}
