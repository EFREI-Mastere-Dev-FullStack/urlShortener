package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type URL struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Original      string             `bson:"original"`
	ShortenedSlug string             `bson:"shortened_slug"`
	Alias         string             `bson:"alias,omitempty"`
	ExpiredAt     time.Time          `bson:"expired_at,omitempty" time_format:"2006-01-02"`
	CreatedAt     time.Time          `bson:"created_at"`
}
