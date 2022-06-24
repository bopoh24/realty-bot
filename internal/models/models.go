package models

import "time"

type User struct {
	ChatID   int64  `json:"chat_id"`
	UserName string `json:"user_name"`
	Name     string `json:"name"`
}

type Ad struct {
	Title    string    `json:"title"`
	Price    int       `json:"price"`
	Link     string    `json:"link"`
	Datetime time.Time `json:"datetime"`
	Location string    `json:"location"`
}
