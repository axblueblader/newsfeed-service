package domains

import "time"

type PostCreateRequest struct {
	Caption  string `json:"caption"`
	ImageUrl string `json:"image_url"`
}

type PostWithComments struct {
	ID        uint      `json:"id"`
	Caption   string    `json:"caption"`
	ImageUrl  string    `json:"image_url"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
	Comments  []Comment `json:"comments"`
}
