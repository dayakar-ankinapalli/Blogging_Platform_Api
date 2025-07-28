package handler

import (
	"encoding/json"
	// "errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gemini/go-blog-api/internal/database"
	"github.com/gemini/go-blog-api/internal/model"
)

// PostHandler handles HTTP requests for blog posts.
type PostHandler struct {
	Store database.Store
}

// NewPostHandler creates a new PostHandler.
func NewPostHandler(s database.Store) *PostHandler {
	return &PostHandler{Store: s}
}

// ServeHTTP routes the request to the appropriate handler method.
func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/posts/")

	// Route to specific handlers based on method and path
	if idStr == "" { // Path is /posts
		switch r.Method {
		case http.MethodGet:
			h.GetAllPosts(w, r)
		case http.MethodPost:
			h.CreatePost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else { // Path is /posts/{id}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			h.GetPost(w, r, id)
		case http.MethodPut:
			h.UpdatePost(w, r, id)
		case http.MethodDelete:
			h.DeletePost(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// CreatePost handles POST /posts
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if post.Title == "" || post.Content == "" {
		http.Error(w, `{"error": "title and content are required"}`, http.StatusBadRequest)
		return
	}

	id, err := h.Store.CreatePost(&post)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// Retrieve the created post to get all fields (like CreatedAt, etc.)
	createdPost, err := h.Store.GetPost(id)
	if err != nil {
		http.Error(w, "Failed to retrieve created post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdPost)
}

// GetAllPosts handles GET /posts
func (h *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")

	posts, err := h.Store.GetAllPosts(term)
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

// GetPost handles GET /posts/{id}
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request, id int64) {
	post, err := h.Store.GetPost(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Post with id %d not found", id), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

// UpdatePost handles PUT /posts/{id}
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request, id int64) {
	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if post.Title == "" || post.Content == "" {
		http.Error(w, `{"error": "title and content are required"}`, http.StatusBadRequest)
		return
	}

	updatedPost, err := h.Store.UpdatePost(id, &post)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update post", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedPost)
}

// DeletePost handles DELETE /posts/{id}
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request, id int64) {
	err := h.Store.DeletePost(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HealthCheckHandler provides a simple health check endpoint.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A simple health check which returns status 200
	data := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
