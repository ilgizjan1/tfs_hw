package ads

type Ad struct {
	ID        int64
	Title     string `json:"title"`
	Text      string `json:"text"`
	AuthorID  int64
	Published bool
}
