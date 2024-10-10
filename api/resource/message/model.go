package message

type Message struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
	Body string `json:"body" binding:"required"`
}
