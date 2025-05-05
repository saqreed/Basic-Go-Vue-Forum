package database

import (
	"forum-go-vue/internal/models"
)

func SaveMessage(message models.ChatMessage) error {
	query := `INSERT INTO chat_messages (content, user_id, reply_to_id, created_at)
			  VALUES ($1, $2, $3, $4) RETURNING id`
	return DB.QueryRow(query, message.Content, message.UserID, message.ReplyToID, message.CreatedAt).Scan(&message.ID)
}

func GetLastMessages(limit int) ([]models.ChatMessage, error) {
	query := `SELECT id, content, user_id, reply_to_id, created_at
			  FROM chat_messages
			  ORDER BY created_at DESC
			  LIMIT $1`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var msg models.ChatMessage
		err := rows.Scan(&msg.ID, &msg.Content, &msg.UserID, &msg.ReplyToID, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	// Reverse the slice to get messages in chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func GetMessageByID(id int) (models.ChatMessage, error) {
	var msg models.ChatMessage
	query := `SELECT id, content, user_id, reply_to_id, created_at
			  FROM chat_messages
			  WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(&msg.ID, &msg.Content, &msg.UserID, &msg.ReplyToID, &msg.CreatedAt)
	return msg, err
}

func GetUsernameByID(id int) (string, error) {
	var username string
	query := `SELECT username FROM users WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(&username)
	return username, err
}
