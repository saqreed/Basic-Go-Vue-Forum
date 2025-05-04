package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

type ProfileHandler struct {
	db *sqlx.DB
}

func NewProfileHandler(db *sqlx.DB) *ProfileHandler {
	return &ProfileHandler{db: db}
}

func (h *ProfileHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	var user struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	err := h.db.QueryRowx("SELECT id, username, email FROM users WHERE id = $1", userID).StructScan(&user)
	if err != nil {
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *ProfileHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	var stats struct {
		Posts    int `json:"posts"`
		Comments int `json:"comments"`
	}

	err := h.db.QueryRowx("SELECT COUNT(*) FROM posts WHERE author_id = $1", userID).Scan(&stats.Posts)
	if err != nil {
		http.Error(w, "Failed to get posts count", http.StatusInternalServerError)
		return
	}

	err = h.db.QueryRowx("SELECT COUNT(*) FROM comments WHERE author_id = $1", userID).Scan(&stats.Comments)
	if err != nil {
		http.Error(w, "Failed to get comments count", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *ProfileHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	var data struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var hashedPassword string
	err := h.db.QueryRowx("SELECT password FROM users WHERE id = $1", userID).Scan(&hashedPassword)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !checkPasswordHash(data.CurrentPassword, hashedPassword) {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	newHashedPassword, err := hashPassword(data.NewPassword)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec("UPDATE users SET password = $1 WHERE id = $2", newHashedPassword, userID)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
