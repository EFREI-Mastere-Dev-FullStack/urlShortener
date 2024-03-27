package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"net/http"
	"os"
	"urlShortener/database"
	"urlShortener/model"
)

const charBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

var BaseURL = os.Getenv("BASE_URL")

func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func ShortenURL(c *gin.Context) {
	var url model.URL

	if err := c.ShouldBind(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Interdire les dates d'expiration dans le passé [url.ExpirationDate]
	// TODO: Interdire le doublon d'alias
	// TODO: Amelioration: Verifier l'unicité du slug en base de données
	url.ShortenedSlug = generateSlug(6)

	db := database.Connection
	_, err := db.Database.Collection("urls").InsertOne(context.Background(), url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "shorten.html", gin.H{"BaseURL": BaseURL, "url": url})
}

func RedirectURL(c *gin.Context) {
	slug := c.Param("slug")

	var url model.URL
	db := database.Connection

	filter := bson.M{"$or": []bson.M{{"shortened_slug": slug}, {"alias": slug}}}
	err := db.Database.Collection("urls").FindOne(context.Background(), filter).Decode(&url)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// TODO: Mettre en place le compteur de clics
	c.Redirect(http.StatusMovedPermanently, url.Original)

}

// private function
func generateSlug(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charBytes[rand.Intn(len(charBytes))]
	}
	return string(b)
}
