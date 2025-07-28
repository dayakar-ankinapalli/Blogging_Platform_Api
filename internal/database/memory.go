package database

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gemini/go-blog-api/internal/model"
)

// MemoryStore is an in-memory implementation of the Store interface.
type MemoryStore struct {
	mu     sync.RWMutex
	posts  map[int64]*model.Post
	nextID int64
}

// NewMemoryStore creates and returns a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		posts:  make(map[int64]*model.Post),
		nextID: 1,
	}
}

// CreatePost adds a new post to the store.
func (s *MemoryStore) CreatePost(post *model.Post) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post.ID = s.nextID
	post.CreatedAt = time.Now().UTC()
	post.UpdatedAt = time.Now().UTC()

	s.posts[post.ID] = post
	s.nextID++

	return post.ID, nil
}

// GetPost retrieves a post by its ID.
func (s *MemoryStore) GetPost(id int64) (*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[id]
	if !ok {
		return nil, fmt.Errorf("post with id %d not found", id)
	}
	return post, nil
}

// GetAllPosts retrieves all posts, with an optional search term filter.
func (s *MemoryStore) GetAllPosts(term string) ([]*model.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	posts := make([]*model.Post, 0, len(s.posts))
	lowerTerm := strings.ToLower(term)

	for _, post := range s.posts {
		if term == "" ||
			strings.Contains(strings.ToLower(post.Title), lowerTerm) ||
			strings.Contains(strings.ToLower(post.Content), lowerTerm) ||
			strings.Contains(strings.ToLower(post.Category), lowerTerm) {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

// UpdatePost updates an existing post.
func (s *MemoryStore) UpdatePost(id int64, post *model.Post) (*model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingPost, ok := s.posts[id]
	if !ok {
		return nil, fmt.Errorf("post with id %d not found", id)
	}

	// Update fields
	existingPost.Title = post.Title
	existingPost.Content = post.Content
	existingPost.Category = post.Category
	existingPost.Tags = post.Tags
	existingPost.UpdatedAt = time.Now().UTC()

	s.posts[id] = existingPost

	return existingPost, nil
}

// DeletePost removes a post from the store.
func (s *MemoryStore) DeletePost(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.posts[id]
	if !ok {
		return fmt.Errorf("post with id %d not found", id)
	}

	delete(s.posts, id)
	return nil
}
