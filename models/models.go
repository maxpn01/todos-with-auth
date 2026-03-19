package models

type Todo struct {
	ID          int64  `json:"id"`
	Text        string `json:"text"`
	IsCompleted bool   `json:"isCompleted"`
}
