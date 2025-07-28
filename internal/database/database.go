package database

import "github.com/gemini/go-blog-api/internal/model"

// Store defines the interface for database operations.
type Store interface {
	CreatePost(post *model.Post) (int64, error)
	GetPost(id int64) (*model.Post, error)
	GetAllPosts(term string) ([]*model.Post, error)
	UpdatePost(id int64, post *model.Post) (*model.Post, error)
	DeletePost(id int64) error
}
