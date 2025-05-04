package middleware

import (
	"context"
	"fmt"
	"forum/config"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var cfg = config.LoadConfig()

type Claims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Auth error: No Authorization header")
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("Auth error: Invalid Authorization header format")
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		log.Printf("Received token: %s", tokenString)
		log.Printf("Using JWT secret: %s", cfg.JWTSecret)

		claims := &struct {
			UserID int64  `json:"user_id"`
			Role   string `json:"role"`
			jwt.StandardClaims
		}{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Auth error: Unexpected signing method: %v", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			log.Printf("Auth error: Token validation failed: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			log.Println("Auth error: Token is invalid")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Printf("Auth success: User ID %d, Role %s", claims.UserID, claims.Role)

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims.Role != "admin" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
