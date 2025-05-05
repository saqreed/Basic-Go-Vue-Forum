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

var clients = make(map[*websocket.Conn]int)
var broadcast = make(chan models.ChatMessageResponse)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	clients[conn] = userID

	// Send last 50 messages
	messages, err := database.GetLastMessages(50)
	if err != nil {
		log.Printf("Error getting messages: %v", err)
		return
	}

	// Get usernames for messages
	for i := range messages {
		username, err := database.GetUsernameByID(messages[i].UserID)
		if err != nil {
			log.Printf("Error getting username: %v", err)
			continue
		}
		messages[i].Username = username

		// If message is a reply, get reply information
		if messages[i].ReplyToID != nil {
			replyMsg, err := database.GetMessageByID(*messages[i].ReplyToID)
			if err != nil {
				log.Printf("Error getting reply message: %v", err)
				continue
			}
			replyUsername, err := database.GetUsernameByID(replyMsg.UserID)
			if err != nil {
				log.Printf("Error getting reply username: %v", err)
				continue
			}
			messages[i].ReplyTo = &models.ReplyTo{
				ID:       replyMsg.ID,
				Content:  replyMsg.Content,
				Username: replyUsername,
			}
		}
	}

	if err := conn.WriteJSON(messages); err != nil {
		log.Printf("Error sending messages: %v", err)
		return
	}

	for {
		var msg models.ChatMessageRequest
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			delete(clients, conn)
			break
		}

		message := models.ChatMessage{
			Content:   msg.Content,
			UserID:    userID,
			ReplyToID: msg.ReplyToID,
			CreatedAt: time.Now(),
		}

		// Save message to database
		if err := database.SaveMessage(message); err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}

		// Get username for the message
		username, err := database.GetUsernameByID(userID)
		if err != nil {
			log.Printf("Error getting username: %v", err)
			continue
		}
		message.Username = username

		// If message is a reply, get reply information
		if message.ReplyToID != nil {
			replyMsg, err := database.GetMessageByID(*message.ReplyToID)
			if err != nil {
				log.Printf("Error getting reply message: %v", err)
				continue
			}
			replyUsername, err := database.GetUsernameByID(replyMsg.UserID)
			if err != nil {
				log.Printf("Error getting reply username: %v", err)
				continue
			}
			message.ReplyTo = &models.ReplyTo{
				ID:       replyMsg.ID,
				Content:  replyMsg.Content,
				Username: replyUsername,
			}
		}

		// Broadcast message to all clients
		for client := range clients {
			if err := client.WriteJSON(message); err != nil {
				log.Printf("Error broadcasting message: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
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
