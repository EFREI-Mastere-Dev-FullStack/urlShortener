package model

import "time"

type URL struct {
	ID            int
	Original      string
	ShortenedSlug string
	Alias         string
	ExpiredAt     time.Time
	CreatedAt     time.Time
}
