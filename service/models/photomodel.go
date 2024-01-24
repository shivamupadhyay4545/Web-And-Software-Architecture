package models

import "time"

type Photo struct {
	Photobytes []byte
	Username   string
	Likes      int32
	NoComments int32
	CreatedAt  time.Time
	PhotoId    string
	Liked      int
}

type Comment struct {
	Content string `json:"content"`
}
