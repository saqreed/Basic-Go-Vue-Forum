package handlers

import (
	"database/sql"
	"encoding/json"
	"forum/internal/database"
	"forum/internal/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int64)
	post.AuthorID = userID
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	query := `INSERT INTO posts (title, content, author_id, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := database.DB.QueryRow(query, post.Title, post.Content, post.AuthorID, post.CreatedAt, post.UpdatedAt).Scan(&post.ID)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	query := `SELECT p.id, p.title, p.content, p.author_id, p.created_at, p.updated_at,
			  u.id, u.username, u.email, u.role, u.created_at,
			  (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) as comments_count
			  FROM posts p
			  JOIN users u ON p.author_id = u.id
			  ORDER BY p.created_at DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.PostResponse
	for rows.Next() {
		var post models.PostResponse
		var author models.UserResponse
		var commentsCount int
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt,
			&author.ID, &author.Username, &author.Email, &author.Role, &author.CreatedAt,
			&commentsCount,
		)
		if err != nil {
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		post.Author = author
		post.CommentsCount = commentsCount
		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	query := `SELECT p.id, p.title, p.content, p.author_id, p.created_at, p.updated_at,
			  u.id, u.username, u.email, u.role, u.created_at
			  FROM posts p
			  JOIN users u ON p.author_id = u.id
			  WHERE p.id = $1`

	var post models.PostResponse
	var author models.UserResponse
	err = database.DB.QueryRow(query, postID).Scan(
		&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt,
		&author.ID, &author.Username, &author.Email, &author.Role, &author.CreatedAt,
	)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	post.Author = author
	json.NewEncoder(w).Encode(post)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int64)
	post.UpdatedAt = time.Now()

	query := `UPDATE posts SET title = $1, content = $2, updated_at = $3
			  WHERE id = $4 AND author_id = $5
			  RETURNING id, title, content, author_id, created_at, updated_at`

	err = database.DB.QueryRow(query, post.Title, post.Content, post.UpdatedAt, postID, userID).Scan(
		&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int64)
	userRole := r.Context().Value("role").(string)

	var query string
	var result sql.Result
	if userRole == "admin" {
		query = `DELETE FROM posts WHERE id = $1`
		result, err = database.DB.Exec(query, postID)
	} else {
		query = `DELETE FROM posts WHERE id = $1 AND author_id = $2`
		result, err = database.DB.Exec(query, postID, userID)
	}

	if err != nil {
		log.Printf("Error deleting post: %v", err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		http.Error(w, "Failed to check deletion result", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Post not found or unauthorized", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte{})
}
