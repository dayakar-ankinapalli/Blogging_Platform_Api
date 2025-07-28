# Go RESTful Blog API

A simple RESTful API for a personal blogging platform built with Go. This project provides basic CRUD (Create, Read, Update, Delete) operations for blog posts.

## Features

- Create, Read, Update, and Delete blog posts.
- List all blog posts.
- Filter blog posts by a search term (searches `title`, `content`, and `category`).
- In-memory data store (no external database required).
- Containerized with Docker.

## Prerequisites

- Go (version 1.21 or later)
- Docker (optional, for containerized deployment)
- `curl` or a REST client like Postman for testing the API.

## Getting Started

### Running Locally

1.  **Clone the repository:**
    ```sh
    git clone <repository-url>
    cd go-blog-api
    ```

2.  **Run the application:**
    You can use the `go run` command or the provided `Makefile`.

    *Using `go run`*:
    ```sh
    go run ./cmd/api
    ```

    *Using `Makefile`*:
    ```sh
    make run
    ```

The API server will start on `http://localhost:8080`.

### Running with Docker

1.  **Build the Docker image:**
    ```sh
    make docker-build
    ```
    This will build an image named `go-blog-api`.

2.  **Run the Docker container:**
    ```sh
    make docker-run
    ```

The API server will be accessible at `http://localhost:8080`.

## API Endpoints

### Post Model

```json
{
  "id": 1,
  "title": "My First Blog Post",
  "content": "This is the content of my first blog post.",
  "category": "Technology",
  "tags": ["Tech", "Programming"],
  "createdAt": "2023-10-27T10:00:00Z",
  "updatedAt": "2023-10-27T10:00:00Z"
}
```

---

### 1. Create a Blog Post

- **Endpoint:** `POST /posts`
- **Description:** Creates a new blog post.
- **Request Body:**
  ```json
  {
    "title": "My First Blog Post",
    "content": "This is the content of my first blog post.",
    "category": "Technology",
    "tags": ["Tech", "Programming"]
  }
  ```
- **Success Response:** `201 Created` with the new post object.
- **Error Response:** `400 Bad Request` for invalid input.

### 2. Get All Blog Posts

- **Endpoint:** `GET /posts`
- **Description:** Retrieves all blog posts. Can be filtered by a search term.
- **Query Parameter:** `term` (optional) - e.g., `GET /posts?term=tech`
- **Success Response:** `200 OK` with an array of post objects.

### 3. Get a Single Blog Post

- **Endpoint:** `GET /posts/{id}`
- **Description:** Retrieves a single blog post by its ID.
- **Success Response:** `200 OK` with the post object.
- **Error Response:** `404 Not Found` if the post does not exist.

### 4. Update a Blog Post

- **Endpoint:** `PUT /posts/{id}`
- **Description:** Updates an existing blog post.
- **Request Body:** Same as the create request.
- **Success Response:** `200 OK` with the updated post object.
- **Error Response:** `404 Not Found` if the post does not exist, `400 Bad Request` for invalid input.

### 5. Delete a Blog Post

- **Endpoint:** `DELETE /posts/{id}`
- **Description:** Deletes a blog post by its ID.
- **Success Response:** `204 No Content`.
- **Error Response:** `404 Not Found` if the post does not exist.

---