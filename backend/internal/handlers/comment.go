package handlers

import (
	"encoding/json"
	"forum/internal/database"
	"forum/internal/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseInt(vars["post_id"], 10, 64)
	if err != nil {
		log.Printf("Invalid post ID: %v", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Проверяем существование поста
	var postExists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)", postID).Scan(&postExists)
	if err != nil {
		log.Printf("Failed to check post existence: %v", err)
		http.Error(w, "Failed to check post existence", http.StatusInternalServerError)
		return
	}

	if !postExists {
		log.Printf("Post with ID %d does not exist", postID)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		log.Printf("User ID not found in context")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	comment.PostID = postID
	comment.AuthorID = userID
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	log.Printf("Creating comment: PostID=%d, AuthorID=%d, Content=%s", comment.PostID, comment.AuthorID, comment.Content)

	query := `INSERT INTO comments (content, post_id, author_id, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err = database.DB.QueryRow(query, comment.Content, comment.PostID, comment.AuthorID, comment.CreatedAt, comment.UpdatedAt).Scan(&comment.ID)
	if err != nil {
		log.Printf("Failed to create comment: %v", err)
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully created comment with ID: %d", comment.ID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseInt(vars["post_id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	query := `SELECT c.id, c.content, c.post_id, c.author_id, c.created_at, c.updated_at,
			  u.id, u.username, u.email, u.role, u.created_at
			  FROM comments c
			  JOIN users u ON c.author_id = u.id
			  WHERE c.post_id = $1
			  ORDER BY c.created_at DESC`

	rows, err := database.DB.Query(query, postID)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []models.CommentResponse
	for rows.Next() {
		var comment models.CommentResponse
		var author models.UserResponse
		var authorID int64
		err := rows.Scan(
			&comment.ID, &comment.Content, &comment.PostID, &authorID, &comment.CreatedAt, &comment.UpdatedAt,
			&author.ID, &author.Username, &author.Email, &author.Role, &author.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Failed to scan comment", http.StatusInternalServerError)
			return
		}
		comment.Author = author
		comments = append(comments, comment)
	}

	json.NewEncoder(w).Encode(comments)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int64)
	comment.UpdatedAt = time.Now()

	query := `UPDATE comments SET content = $1, updated_at = $2
			  WHERE id = $3 AND author_id = $4
			  RETURNING id, content, post_id, author_id, created_at, updated_at`

	err = database.DB.QueryRow(query, comment.Content, comment.UpdatedAt, commentID, userID).Scan(
		&comment.ID, &comment.Content, &comment.PostID, &comment.AuthorID, &comment.CreatedAt, &comment.UpdatedAt,
	)
	if err != nil {
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comment)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int64)
	userRole := r.Context().Value("role").(string)

	var query string
	if userRole == "admin" {
		query = `DELETE FROM comments WHERE id = $1`
		err = database.DB.QueryRow(query, commentID).Scan()
	} else {
		query = `DELETE FROM comments WHERE id = $1 AND author_id = $2`
		err = database.DB.QueryRow(query, commentID, userID).Scan()
	}

	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
