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

type PostGetAllRequest struct {
	CursorID     *uint `form:"cursor_id"`
	CommentCount *int  `form:"comment_count"`
	PageSize     int   `form:"page_size"`
}

type PostsPagedResult struct {
	Posts        []PostWithComments `json:"posts"`
	NextCursor   *uint              `json:"next_cursor"`
	CommentCount *int               `json:"comment_count"`
	PageSize     int                `json:"page_size"`
}
