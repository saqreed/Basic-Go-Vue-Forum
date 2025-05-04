package handlers

import (
	"forum/internal/database"
	"forum/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Printf("WebSocket connection attempt from origin: %s", origin)
		return true // В продакшене нужно настроить правильную проверку origin
	},
	HandshakeTimeout: 10 * time.Second,
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.ChatMessageResponse)

func HandleChat(w http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming WebSocket request from %s", r.RemoteAddr)
	log.Printf("Request headers: %v", r.Header)

	// Проверяем, что это WebSocket запрос
	if !websocket.IsWebSocketUpgrade(r) {
		log.Printf("Not a WebSocket upgrade request")
		http.Error(w, "Not a websocket handshake", http.StatusBadRequest)
		return
	}

	// Получаем user_id из контекста перед установкой соединения
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		log.Printf("User not authenticated")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer func() {
		conn.Close()
		delete(clients, conn)
		log.Printf("Client disconnected. Total clients: %d", len(clients))
	}()

	clients[conn] = true
	log.Printf("New client connected. Total clients: %d", len(clients))

	// Отправляем последние 50 сообщений новому клиенту
	messages, err := getLastMessages(50)
	if err != nil {
		log.Printf("Failed to get last messages: %v", err)
	} else {
		for _, msg := range messages {
			if err := conn.WriteJSON(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
				return
			}
		}
	}

	// Устанавливаем таймауты для чтения и записи
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Запускаем пинг-понг
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("Failed to send ping: %v", err)
					return
				}
			}
		}
	}()

	// Обработка сообщений
	for {
		var msg models.ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			return
		}

		msg.UserID = userID
		msg.CreatedAt = time.Now()

		// Сохраняем сообщение в базу данных
		query := `INSERT INTO chat_messages (content, user_id, created_at)
				  VALUES ($1, $2, $3) RETURNING id`
		err = database.DB.QueryRow(query, msg.Content, msg.UserID, msg.CreatedAt).Scan(&msg.ID)
		if err != nil {
			log.Printf("Failed to save message: %v", err)
			continue
		}

		// Получаем информацию о пользователе
		response := models.ChatMessageResponse{
			ID:        msg.ID,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		}

		query = `SELECT id, username, email, role, created_at
				 FROM users WHERE id = $1`
		var user models.UserResponse
		err = database.DB.QueryRow(query, msg.UserID).Scan(
			&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt,
		)
		if err != nil {
			log.Printf("Failed to get user info: %v", err)
			continue
		}
		response.User = user

		// Отправляем сообщение всем клиентам
		broadcast <- response
	}
}

func getLastMessages(limit int) ([]models.ChatMessageResponse, error) {
	query := `SELECT m.id, m.content, m.created_at,
			  u.id, u.username, u.email, u.role, u.created_at
			  FROM chat_messages m
			  JOIN users u ON m.user_id = u.id
			  ORDER BY m.created_at DESC
			  LIMIT $1`

	rows, err := database.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.ChatMessageResponse
	for rows.Next() {
		var msg models.ChatMessageResponse
		var user models.UserResponse
		err := rows.Scan(
			&msg.ID, &msg.Content, &msg.CreatedAt,
			&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		msg.User = user
		messages = append(messages, msg)
	}

	return messages, nil
}

func BroadcastMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error broadcasting message: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
