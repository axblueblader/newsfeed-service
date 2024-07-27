package domains

import "time"

type CreateCommentRequest struct {
	Content string `json:"content"`
}

type Comment struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
}
