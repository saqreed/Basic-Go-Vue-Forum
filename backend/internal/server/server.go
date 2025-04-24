package server

import (
	"forum/internal/auth"
	"forum/internal/handlers"
	"forum/internal/middleware"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	Router *mux.Router
	DB     *sqlx.DB
}

func NewServer(db *sqlx.DB) *Server {
	s := &Server{
		Router: mux.NewRouter(),
		DB:     db,
	}
	s.SetupRoutes()
	return s
}

func (s *Server) SetupRoutes() {
	router := s.Router

	// Auth service
	authService := auth.NewAuthService(s.DB)
	router.HandleFunc("/api/register", authService.Register).Methods("POST")
	router.HandleFunc("/api/login", authService.Login).Methods("POST")

	// Profile routes
	profileHandler := handlers.NewProfileHandler(s.DB)
	router.Handle("/api/profile", middleware.AuthMiddleware(http.HandlerFunc(profileHandler.GetUserProfile))).Methods("GET")
	router.Handle("/api/profile/stats", middleware.AuthMiddleware(http.HandlerFunc(profileHandler.GetUserStats))).Methods("GET")
	router.Handle("/api/profile/password", middleware.AuthMiddleware(http.HandlerFunc(profileHandler.ChangePassword))).Methods("PUT")

	// Posts routes
	router.HandleFunc("/api/posts", handlers.GetPosts).Methods("GET")
	router.HandleFunc("/api/posts/{id}", handlers.GetPost).Methods("GET")
	router.Handle("/api/posts", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreatePost))).Methods("POST")
	router.Handle("/api/posts/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdatePost))).Methods("PUT")
	router.Handle("/api/posts/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeletePost))).Methods("DELETE")

	// Comments routes
	router.HandleFunc("/api/posts/{post_id}/comments", handlers.GetComments).Methods("GET")
	router.Handle("/api/posts/{post_id}/comments", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreateComment))).Methods("POST")
	router.Handle("/api/comments/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateComment))).Methods("PUT")
	router.Handle("/api/comments/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteComment))).Methods("DELETE")

	// Admin routes
	adminRouter := router.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(middleware.AdminMiddleware)
	adminRouter.HandleFunc("/users", handlers.GetAllUsers).Methods("GET")
	adminRouter.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	adminRouter.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	adminRouter.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")
	adminRouter.HandleFunc("/posts", handlers.GetAllPosts).Methods("GET")
	adminRouter.HandleFunc("/comments", handlers.GetAllComments).Methods("GET")
}
