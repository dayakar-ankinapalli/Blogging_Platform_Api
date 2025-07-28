package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gemini/go-blog-api/internal/model"
)

// mockStore is a mock implementation of the database.Store for testing purposes.
type mockStore struct {
	posts  map[int64]*model.Post
	nextID int64
	err    error // To simulate database errors
}

func newMockStore() *mockStore {
	return &mockStore{
		posts:  make(map[int64]*model.Post),
		nextID: 1,
	}
}

func (m *mockStore) CreatePost(post *model.Post) (int64, error) {
	if m.err != nil {
		return 0, m.err
	}
	id := m.nextID
	post.ID = id
	now := time.Now().UTC()
	post.CreatedAt = now
	post.UpdatedAt = now
	m.posts[id] = post
	m.nextID++
	return id, nil
}

func (m *mockStore) GetPost(id int64) (*model.Post, error) {
	if m.err != nil {
		return nil, m.err
	}
	post, ok := m.posts[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return post, nil
}

func (m *mockStore) GetAllPosts(term string) ([]*model.Post, error) {
	if m.err != nil {
		return nil, m.err
	}
	posts := make([]*model.Post, 0, len(m.posts))
	for _, p := range m.posts {
		posts = append(posts, p)
	}
	return posts, nil
}

func (m *mockStore) UpdatePost(id int64, post *model.Post) (*model.Post, error) {
	if m.err != nil {
		return nil, m.err
	}
	_, ok := m.posts[id]
	if !ok {
		return nil, errors.New("not found")
	}
	post.ID = id
	post.UpdatedAt = time.Now().UTC()
	m.posts[id] = post
	return post, nil
}

func (m *mockStore) DeletePost(id int64) error {
	if m.err != nil {
		return m.err
	}
	_, ok := m.posts[id]
	if !ok {
		return errors.New("not found")
	}
	delete(m.posts, id)
	return nil
}

func TestPostHandler(t *testing.T) {
	store := newMockStore()
	handler := NewPostHandler(store)

	// Seed a post for GET, UPDATE, DELETE tests
	initialPost := &model.Post{
		Title:   "Initial Post",
		Content: "Initial Content",
	}
	store.CreatePost(initialPost)

	t.Run("CreatePost", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			postData := map[string]interface{}{
				"title":   "New Post",
				"content": "New Content",
			}
			body, _ := json.Marshal(postData)

			req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusCreated {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
			}

			var createdPost model.Post
			json.Unmarshal(rr.Body.Bytes(), &createdPost)
			if createdPost.Title != "New Post" {
				t.Errorf("handler returned unexpected body: got title %v want %v", createdPost.Title, "New Post")
			}
		})

		t.Run("bad request - invalid json", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader([]byte("{invalid")))
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
			}
		})

		t.Run("bad request - missing title", func(t *testing.T) {
			postData := map[string]interface{}{"content": "Some content"}
			body, _ := json.Marshal(postData)
			req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(body))
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
			}
		})
	})

	t.Run("GetPost", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/posts/1", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}

			var post model.Post
			json.Unmarshal(rr.Body.Bytes(), &post)
			if post.ID != 1 {
				t.Errorf("handler returned wrong post ID: got %v want %v", post.ID, 1)
			}
		})

		t.Run("not found", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/posts/999", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusNotFound {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
			}
		})
	})

	t.Run("GetAllPosts", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var posts []model.Post
		json.Unmarshal(rr.Body.Bytes(), &posts)
		if len(posts) == 0 {
			t.Error("handler returned no posts, expected at least one")
		}
	})

	t.Run("UpdatePost", func(t *testing.T) {
		updateData := map[string]interface{}{
			"title":   "Updated Title",
			"content": "Updated Content",
		}
		body, _ := json.Marshal(updateData)

		req := httptest.NewRequest(http.MethodPut, "/posts/1", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var updatedPost model.Post
		json.Unmarshal(rr.Body.Bytes(), &updatedPost)
		if updatedPost.Title != "Updated Title" {
			t.Errorf("handler did not update title: got %v want %v", updatedPost.Title, "Updated Title")
		}
	})

	t.Run("DeletePost", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			// Use a new post ID to avoid interfering with other tests
			postToDelete := &model.Post{Title: "To Delete", Content: "Content"}
			id, _ := store.CreatePost(postToDelete)

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/posts/%d", id), nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusNoContent {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
			}
		})

		t.Run("not found", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/posts/999", nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusNotFound {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
			}
		})
	})
}