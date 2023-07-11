package models

import "time"

type Ad struct {
	ID           int64
	Title        string `json:"title"`
	Text         string `json:"text"`
	UserID       int64
	Published    bool
	DateCreation time.Time
	DateUpdate   time.Time
}
