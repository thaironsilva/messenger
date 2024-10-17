package message

import "time"

type Message struct {
	Id         string
	SenderId   string    `json:"senderId" binding:"required"`
	ReceiverId string    `json:"receiverId" binding:"required"`
	Body       string    `json:"body" binding:"required"`
	CreatedAt  time.Time `json:"createdAt" binding:"required"`
}
