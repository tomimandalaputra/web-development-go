package main

import (
	"database/sql"
	"errors"
	"math"
	"strings"
	"time"
)

var (
	ErrDuplicatePostTitle = errors.New("Duplicate post title ")
	ErrDuplicateVote      = errors.New("Duplicate vote ")
)

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

func (f *Filter) Validate() error {
	if f.Page <= 0 || f.PageSize >= 10_000_000 {
		return errors.New("Invalid page range: 1 to 10 million")
	}

	if f.PageSize <= 0 || f.PageSize >= 100 {
		return errors.New("Invalid page range: 1 to 100 max")
	}

	return nil
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

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	meta := Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}

	meta.NextPage = meta.CurrentPage + 1
	meta.PrevPage = meta.CurrentPage - 1

	if meta.CurrentPage <= meta.FirstPage {
		meta.PrevPage = 0
	}
	if meta.CurrentPage >= meta.NextPage {
		meta.NextPage = 0
	}

	return meta
}

type PostRepository interface {
	CreatePost(title, url string, userID int) (int, error)
	AddComment(userID, postID int, body string) (int, error)
	AddVote(userID, postID int) error
	GetAll(filter Filter) ([]Post, Metadata, error)
	GetByID(id int) (*Post, error)
	GetComments(postID int) ([]Comment, error)
}

type SQLPostRepository struct {
	db *sql.DB
}

// NewSQLPostRepository create new SQLPostRepository type
func NewSQLPostRepository(db *sql.DB) *SQLPostRepository {
	return &SQLPostRepository{
		db: db,
	}
}

func (r *SQLPostRepository) CreatePost(title, url string, userID int) (int, error) {
	stmt := "INSERT INTO posts (title, url, user_id) VALUES (?,?,?)"

	result, err := r.db.Exec(stmt, title, url, userID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: post.title") {
			return 0, ErrDuplicatePostTitle
		}
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(postID), nil
}

func (r *SQLPostRepository) AddComment(userID, postID int, body string) (int, error) {
	stmt := "INSERT INTO comments (user_id, post_id, body) VALUES (?,?,?)"

	result, err := r.db.Exec(stmt, userID, postID, body)
	if err != nil {
		return 0, err
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(commentID), nil
}

func (r *SQLPostRepository) AddVote(userID, postID int) error {
	stmt := "INSERT INTO votes (user_id, post_id) VALUES (?, ?)"

	_, err := r.db.Exec(stmt, userID, postID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "PRIMARY KEY constraint failed") {
			return ErrDuplicateVote
		}

		return err
	}

	return nil
}

func (r *SQLPostRepository) GetByID(id int) (*Post, error) {
	query := `
		SELECT p.id, p.title, p.url, p.user_id, p.created_at,
		u.name as user_name,
		COUNT(DISTINCT c.id) AS comment_count,
		COUNT(DISTINCT v.user_id) AS voute_count
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN votes v ON p.id = v.post_id
		WHERE p.id = ?
		GROUP BY p.id, p.title, p.url, p.user_id, p.created_at, u.name
	`

	row := r.db.QueryRow(query, id)

	var post Post
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.URL,
		&post.UserID,
		&post.CreatedAt,
		&post.UserName,
		&post.CommentCount,
		&post.VoteCount,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}
