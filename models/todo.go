package models

type Todo struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userID"`
	Text        string `json:"text"`
	IsCompleted bool   `json:"isCompleted"`
}
