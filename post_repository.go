package main

import "time"

type Post struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	UserID       int       `json:"user_id"`
	UserName     string    `json:"user_name"`
	CommentCount int       `json:"comment_count"`
	VoteCount    int       `json:"vote_count"`
	TotalRecords int       `json:"total_records"`
	CreatedAt    time.Time `json:"created_at"`
}

type Comment struct {
	ID        int       `json:"id"`
	Body      string    `json:"body"`
	UserID    int       `json:"user_id"`
	PostID    int       `json:"post_id"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

type Filter struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	OrderBy  int `json:"order_by"`
	Query    int `json:"query"`
}

type Metadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	NextPage     int `json:"next_page"`
	PrevPage     int `json:"prev_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}

type PostRepository interface {
	CreatePost(title, url string, userID int) (int, error)
	AddComment(userID, postID int, body string) (int, error)
	AddVote(userID, postID int, body string) error
	GetAll(filter Filter) ([]Post, Metadata, error)
	GetByID(id int) (*Post, error)
	GetComments(postID int) ([]Comment, error)
}
