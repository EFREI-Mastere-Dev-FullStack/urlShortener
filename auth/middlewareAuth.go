package auth

import (
	"context"
	"net/http"
	"urlShortener/database"
	"urlShortener/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/subinoybiswas/goenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RequiredAuth(c *gin.Context) {
	conn := database.Connection
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/")
		c.Abort()
		return
	}

	secretKey, _ := goenv.GetEnv("SECRET_KEY")
	if secretKey == "" {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error (missing secret key)"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
			c.Abort()
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["sub"].(string)
		var user model.User
		err := conn.Database.Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if user.ID == primitive.NilObjectID {

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized user"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	} else {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		c.Abort()
		return
	}
}
