package controller

import (
	"net/http"
	"time"
	"urlShortener/model"

	"github.com/gin-gonic/gin"
	"github.com/subinoybiswas/goenv"
)

var BaseURL, _ = goenv.GetEnv("BASE_URL")

func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func ShortenURL(c *gin.Context) {
	var url model.URL
	if err := c.ShouldBind(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := url.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "shorten.html", gin.H{"BaseURL": BaseURL, "url": url})
}

func RedirectURL(c *gin.Context) {
	slug := c.Param("slug")

	var url model.URL
	err := url.DecodeSlugAndCount(slug)

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
