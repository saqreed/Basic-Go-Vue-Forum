package handlers

import (
	"encoding/json"
	"forum/internal/database"
	"forum/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, username, email, role, created_at FROM users ORDER BY created_at DESC`

	rows, err := database.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to scan user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `SELECT id, username, email, role, created_at FROM users WHERE id = $1`

	var user models.UserResponse
	err = database.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE users SET username = $1, email = $2, role = $3
			  WHERE id = $4
			  RETURNING id, username, email, role, created_at`

	err = database.DB.QueryRow(query, user.Username, user.Email, user.Role, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM users WHERE id = $1`
	err = database.DB.QueryRow(query, userID).Scan()
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	query := `SELECT p.id, p.title, p.content, p.author_id, p.created_at, p.updated_at,
			  u.id, u.username, u.email, u.role, u.created_at
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
		err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt,
			&author.ID, &author.Username, &author.Email, &author.Role, &author.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		post.Author = author
		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}

func GetAllComments(w http.ResponseWriter, r *http.Request) {
	query := `SELECT c.id, c.content, c.post_id, c.author_id, c.created_at, c.updated_at,
			  u.id, u.username, u.email, u.role, u.created_at
			  FROM comments c
			  JOIN users u ON c.author_id = u.id
			  ORDER BY c.created_at DESC`

	rows, err := database.DB.Query(query)
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

func GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	query := `SELECT c.id, c.content, c.post_id, c.author_id, c.created_at, c.updated_at,
			  u.id, u.username, u.email, u.role, u.created_at
			  FROM comments c
			  JOIN users u ON c.author_id = u.id
			  WHERE c.id = $1`

	var comment models.CommentResponse
	var author models.UserResponse
	var authorID int64
	err = database.DB.QueryRow(query, commentID).Scan(
		&comment.ID, &comment.Content, &comment.PostID, &authorID, &comment.CreatedAt, &comment.UpdatedAt,
		&author.ID, &author.Username, &author.Email, &author.Role, &author.CreatedAt,
	)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}
	comment.Author = author
	json.NewEncoder(w).Encode(comment)
}
