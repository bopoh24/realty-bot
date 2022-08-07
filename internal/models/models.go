package models

import "time"

type ChatID int64

type Chat struct {
	ChatID   ChatID `json:"chat_id"`
	UserName string `json:"user_name"`
	Name     string `json:"name"`
}

type AdLink string

type Ad struct {
	Title    string    `json:"title"`
	Price    int       `json:"price"`
	Link     AdLink    `json:"link"`
	Datetime time.Time `json:"datetime"`
	Location string    `json:"location"`
}
