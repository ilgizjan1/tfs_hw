package ads

import "time"

type Ad struct {
	ID           int64
	Title        string `json:"title"`
	Text         string `json:"text"`
	AuthorID     int64
	Published    bool
	DateCreation time.Time
	DateUpdate   time.Time
}
