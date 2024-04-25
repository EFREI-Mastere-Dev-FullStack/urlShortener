package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"net/http"
	"os"
	"time"
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

	db := database.Connection

	err := createURL(db, &url)
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
	update := bson.M{"$inc": bson.M{"count": 1}}

	err := db.Database.Collection("urls").FindOneAndUpdate(context.Background(), filter, update).Decode(&url)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	if !url.ExpiredAt.IsZero() && url.ExpiredAt.Before(time.Now().Local()) {
		c.JSON(http.StatusGone, gin.H{"error": "URL expired"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url.Original)
}

// private function
func createURL(db *database.MongoDBConnection, url *model.URL) error {
	if url.Alias != "" {
		filter := bson.M{"$or": []bson.M{{"shortened_slug": url.Alias}, {"alias": url.Alias}}}
		err := db.Database.Collection("urls").FindOne(context.Background(), filter).Err()
		if err == nil {
			return fmt.Errorf("alias already in use")
		}
	}

	url.ShortenedSlug = generateSlug(db, 6)
	url.CreatedAt = time.Now().Local()
	url.Count = 0

	if !url.ExpiredAt.IsZero() && url.ExpiredAt.Before(url.CreatedAt) {
		return fmt.Errorf("expiration date can't be in the past")
	}

	_, err := db.Database.Collection("urls").InsertOne(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error inserting url into database: %v", err)
	}

	return nil
}

func generateSlug(db *database.MongoDBConnection, n int) string {
	var slug string
	for {
		b := make([]byte, n)
		for i := range b {
			b[i] = charBytes[rand.Intn(len(charBytes))]
		}
		slug = string(b)

		filter := bson.M{"shortened_slug": slug}
		err := db.Database.Collection("urls").FindOne(context.Background(), filter).Err()
		if err != nil {
			break
		}
	}

	return slug
}
