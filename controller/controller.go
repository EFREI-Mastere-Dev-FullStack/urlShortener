package controller

import (
	"context"
	"github.com/gin-gonic/gin"
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

	// TODO: Amelioration: Verifier l'unicité du slug en base de données
	url.ShortenedSlug = generateSlug(6)

	// TODO: Interdire les dates d'expiration dans le passé [url.ExpirationDate]

	db := database.Connection
	_, err := db.Database.Collection("urls").InsertOne(context.Background(), url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "shorten.html", gin.H{"BaseURL": BaseURL, "url": url})
}

// private function
func generateSlug(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charBytes[rand.Intn(len(charBytes))]
	}
	return string(b)
}
