package auth

import (
	"context"
	"errors"
	"time"
	"urlShortener/database"
	"urlShortener/model"

	"github.com/golang-jwt/jwt"
	"github.com/subinoybiswas/goenv"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Register(username string, password string) error {
	conn := database.Connection

	if isUserAlreadyExists(username) {
		return errors.New("username already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.User{
		Username: username,
		Password: string(hashedPassword),
	}

	_, err = conn.Database.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}

func Login(username string, password string) (string, error) {

	conn := database.Connection

	var user model.User
	err := conn.Database.Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	var secret, _ = goenv.GetEnv("SECRET_KEY")

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("error generating token")
	}
	return tokenString, nil
}

func isUserAlreadyExists(username string) bool {
	conn := database.Connection

	var existingUser model.User
	err := conn.Database.Collection("users").FindOne(context.Background(), bson.M{"username": username}).Decode(&existingUser)
	if err != nil {
		return false
	}
	return true
}
