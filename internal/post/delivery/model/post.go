package model

import "time"

// [GET].
type PostsView struct {
	PostsView []*PostView `json:"posts"`
}

// [GET].
type PostView struct {
	PostID    string    `json:"post_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Image     string    `json:"image,omitempty"`
	User      string    `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
