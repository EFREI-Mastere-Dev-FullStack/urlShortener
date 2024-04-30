package model

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"urlShortener/database"

	"github.com/subinoybiswas/goenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URL struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Original      string             `bson:"original"`
	ShortenedSlug string             `bson:"shortened_slug"`
	Alias         string             `bson:"alias,omitempty"`
	ExpiredAt     time.Time          `bson:"expired_at,omitempty" time_format:"2006-01-02"`
	CreatedAt     time.Time          `bson:"created_at"`
	Count         int                `bson:"count"`
	UserID        primitive.ObjectID `bson:"user_id,omitempty"`
}

const charBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

func (url *URL) Save() error {
	var db = database.Connection

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

func (url *URL) DecodeSlugAndCount(slug string) error {
	var db = database.Connection

	filter := bson.M{"$or": []bson.M{{"shortened_slug": slug}, {"alias": slug}}}
	update := bson.M{"$inc": bson.M{"count": 1}}

	err := db.Database.Collection("urls").FindOneAndUpdate(context.Background(), filter, update).Decode(&url)

	return err
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

type URLWithShortened struct {
	Original      string    `bson:"original"`
	ShortenedSlug string    `bson:"shortened_slug"`
	Alias         string    `bson:"alias,omitempty"`
	ExpiredAt     time.Time `bson:"expired_at,omitempty" time_format:"2006-01-02"`
	URLShorten    string    `bson:"-"`
	URLWithAlias  string    `bson:"-"`
}

func GetUrlsByUserId(userID primitive.ObjectID) ([]URLWithShortened, error) {
	var db = database.Connection

	var urls []URLWithShortened
	filter := bson.M{"user_id": userID}
	cursor, err := db.Database.Collection("urls").Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("error getting urls from database: %v", err)
	}

	baseURL, err := goenv.GetEnv("BASE_URL")
	if err != nil {
		return nil, fmt.Errorf("error getting base url from environment: %v", err)
	}

	for cursor.Next(context.Background()) {
		var url URL
		if err := cursor.Decode(&url); err != nil {
			return nil, fmt.Errorf("error decoding url from database: %v", err)
		}

		var urlShorten string
		if url.Alias != "" {
			urlShorten = fmt.Sprintf("%s/%s", baseURL, url.Alias)
		} else {
			urlShorten = fmt.Sprintf("%s/%s", baseURL, url.ShortenedSlug)
		}

		urlWithShortened := URLWithShortened{
			Original:      url.Original,
			ShortenedSlug: url.ShortenedSlug,
			Alias:         url.Alias,
			ExpiredAt:     url.ExpiredAt,
			URLShorten:    urlShorten,
		}
		urls = append(urls, urlWithShortened)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cursor: %v", err)
	}

	return urls, nil
}
