package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"urlShortener/database"
	"urlShortener/model"
)

func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func ShortenURL(c *gin.Context) {
	var url model.URL

	if err := c.ShouldBind(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Generer un slug [faudrait faire attention a ce que le slug soit unique en BDD]

	db := database.Connection
	_, err := db.Database.Collection("urls").InsertOne(context.Background(), url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL shorted"})
}
