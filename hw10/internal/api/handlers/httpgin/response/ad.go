package response

type AdResponse struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Text         string `json:"text"`
	UserID       int64  `json:"user_id"`
	Published    bool   `json:"published"`
	DateCreation string `json:"date_creation"`
	DateUpdate   string `json:"date_update"`
}
