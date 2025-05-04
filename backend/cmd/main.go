package main

import (
	"context"
	"forum/internal/auth"
	"forum/internal/database"
	"forum/internal/handlers"
	"forum/internal/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	r := mux.NewRouter()
	authService := auth.NewAuthService(database.DB)
	r.HandleFunc("/api/register", authService.Register).Methods("POST")
	r.HandleFunc("/api/login", authService.Login).Methods("POST")
	r.Handle("/api/change-password", middleware.AuthMiddleware(http.HandlerFunc(authService.ChangePassword))).Methods("POST")
	r.HandleFunc("/api/posts", handlers.GetPosts).Methods("GET")
	r.HandleFunc("/api/posts/{id}", handlers.GetPost).Methods("GET")
	r.HandleFunc("/api/posts/{post_id}/comments", handlers.GetComments).Methods("GET")

	// WebSocket handler with proper middleware order
	wsHandler := middleware.WebSocketAuthMiddleware(http.HandlerFunc(handlers.HandleChat))
	r.Handle("/ws/chat", wsHandler).Methods("GET")

	authRouter := r.PathPrefix("/api").Subrouter()
	authRouter.Use(middleware.AuthMiddleware)

	profileHandler := handlers.NewProfileHandler(database.DB)
	authRouter.HandleFunc("/profile", profileHandler.GetUserProfile).Methods("GET")
	authRouter.HandleFunc("/profile/stats", profileHandler.GetUserStats).Methods("GET")
	authRouter.HandleFunc("/profile/password", profileHandler.ChangePassword).Methods("PUT")

	authRouter.HandleFunc("/posts", handlers.CreatePost).Methods("POST")
	authRouter.HandleFunc("/posts/{id}", handlers.UpdatePost).Methods("PUT")
	authRouter.HandleFunc("/posts/{id}", handlers.DeletePost).Methods("DELETE")
	authRouter.HandleFunc("/posts/{post_id}/comments", handlers.CreateComment).Methods("POST")
	authRouter.HandleFunc("/comments/{id}", handlers.UpdateComment).Methods("PUT")
	authRouter.HandleFunc("/comments/{id}", handlers.DeleteComment).Methods("DELETE")

	adminRouter := r.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(middleware.AdminMiddleware)

	adminRouter.HandleFunc("/users", handlers.GetAllUsers).Methods("GET")
	adminRouter.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	adminRouter.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	adminRouter.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")

	adminRouter.HandleFunc("/posts", handlers.GetAllPosts).Methods("GET")
	adminRouter.HandleFunc("/posts/{id}", handlers.GetPost).Methods("GET")
	adminRouter.HandleFunc("/posts/{id}", handlers.UpdatePost).Methods("PUT")
	adminRouter.HandleFunc("/posts/{id}", handlers.DeletePost).Methods("DELETE")

	adminRouter.HandleFunc("/comments", handlers.GetAllComments).Methods("GET")
	adminRouter.HandleFunc("/comments/{id}", handlers.GetComment).Methods("GET")
	adminRouter.HandleFunc("/comments/{id}", handlers.UpdateComment).Methods("PUT")
	adminRouter.HandleFunc("/comments/{id}", handlers.DeleteComment).Methods("DELETE")

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Для WebSocket соединений
			if r.Header.Get("Upgrade") == "websocket" {
				next.ServeHTTP(w, r)
				return
			}

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	handler := corsMiddleware(r)

	srv := &http.Server{
		Addr:         ":8081",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Server starting on :8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Запускаем обработку широковещательных сообщений
	go handlers.BroadcastMessages()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
