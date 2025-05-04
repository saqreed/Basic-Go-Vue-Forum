package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func WebSocketAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовки CORS до проверки WebSocket
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Пропускаем OPTIONS запросы
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Проверяем, что это WebSocket запрос
		if !strings.Contains(strings.ToLower(r.Header.Get("Upgrade")), "websocket") {
			log.Println("Not a WebSocket request")
			http.Error(w, "Not a WebSocket request", http.StatusBadRequest)
			return
		}

		// Извлекаем токен из RawQuery
		token := ""
		rawQuery := r.URL.RawQuery
		if rawQuery != "" {
			params := strings.Split(rawQuery, "&")
			for _, param := range params {
				if strings.HasPrefix(param, "token=") {
					token = strings.TrimPrefix(param, "token=")
					break
				}
			}
		}

		if token == "" {
			// Если токена нет в URL, пробуем получить из заголовка
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("WebSocket auth error: No token provided")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Println("WebSocket auth error: Invalid Authorization header format")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			token = parts[1]
		}

		claims := &struct {
			UserID int64  `json:"user_id"`
			Role   string `json:"role"`
			jwt.StandardClaims
		}{}

		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("WebSocket auth error: Unexpected signing method: %v", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			log.Printf("WebSocket auth error: Token validation failed: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !parsedToken.Valid {
			log.Println("WebSocket auth error: Token is invalid")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("WebSocket auth success: User ID %d, Role %s", claims.UserID, claims.Role)

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
