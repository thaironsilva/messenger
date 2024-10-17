package message

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAll(sender_id string, receiver_id string) ([]Message, error) {
	query := "SELECT * FROM messages WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1) LIMIT 20"
	rows, err := r.db.Query(query, sender_id, receiver_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message

	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.Id, &message.SenderId, &message.ReceiverId, &message.Body, &message.CreatedAt); err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *Repository) Create(newMessage Message) error {
	query := "INSERT INTO messages (sender_id, receiver_id, body, created_at) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, newMessage.SenderId, newMessage.ReceiverId, newMessage.Body, newMessage.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
