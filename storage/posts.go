package storage

import (
	"gorm.io/gorm"
	"log"
	"newsfeed-service/config"
	"newsfeed-service/domains"
	"newsfeed-service/models"
)

type postDB struct {
	db *gorm.DB
}

func NewPostDB(db *gorm.DB) PostDB {
	return &postDB{db: db}
}

type PostDB interface {
	Create(post *models.Post) (*models.Post, error)
	GetAllWithComments(userID string, cursor domains.PostGetAllRequest) ([]models.Post, *domains.PostGetAllRequest, error)
}

func (dao postDB) Create(post *models.Post) (*models.Post, error) {
	err := dao.db.Create(post).Error
	return post, err
}

type PostResult struct {
	models.Post
	CommentCount int
}

func (dao postDB) GetAllWithComments(userID string, cursor domains.PostGetAllRequest) ([]models.Post, *domains.PostGetAllRequest, error) {
	var posts []models.Post
	var postResults []PostResult
	var postIDs []uint
	// we can get next cursor or no cursor based on len(result) == pageSize
	pageSize := cursor.PageSize
	limit := pageSize + 1

	args := map[string]interface{}{
		"limit":  limit,
		"userID": userID,
	}

	// without the outer select, it will miss posts with no comments
	// get posts sort by highest comments, then sort by id and created_at for pagination
	// cursor needs to be unique and based around comment_count
	selectClause := `
            SELECT p.*, COALESCE(c.comment_count, 0) AS comment_count
            FROM posts p
            LEFT JOIN (
                SELECT post_id, COUNT(*) AS comment_count
                FROM comments
                GROUP BY post_id
            ) c ON p.id = c.post_id
          `
	whereClause := `WHERE p.creator = @userID`
	if cursor.CommentCount != nil && cursor.CursorID != nil {
		whereClause = whereClause + ` AND comment_count < @cmtCount OR (comment_count = @cmtCount AND p.id <= @cursorID) `
		args["cursorID"] = *cursor.CursorID
		args["cmtCount"] = *cursor.CommentCount
	}
	orderAndLimit := ` ORDER BY comment_count DESC, p.created_at DESC, p.id DESC LIMIT @limit`

	rawSql := selectClause + whereClause + orderAndLimit

	err := dao.db.Raw(rawSql, args).Scan(&postResults).Error

	if err != nil {
		return nil, nil, err
	}

	log.Println(postResults)
	if len(postResults) == 0 {
		return posts, nil, nil
	}

	postsMap := make(map[uint]int)
	for i := range postResults {
		result := postResults[i]
		posts = append(posts, result.Post)
		postsMap[result.Post.ID] = i
		postIDs = append(postIDs, result.Post.ID)
	}

	commentsPerPostsSql := `
	SELECT c.*
	FROM (
		SELECT *, ROW_NUMBER() OVER (PARTITION BY post_id ORDER BY created_at DESC) AS row_number
		FROM comments 
		WHERE post_id IN ?
		) as c
	WHERE c.row_number <= ?`

	var comments []models.Comment
	err = dao.db.Raw(commentsPerPostsSql, postIDs, config.Env().CommentsLimit).Scan(&comments).Error
	if err != nil {
		return nil, nil, err
	}
	for i := range comments {
		comment := comments[i]
		post := &posts[postsMap[comment.PostID]]
		if len(post.Comments) == 0 {
			post.Comments = []models.Comment{}
		}
		post.Comments = append(post.Comments, comment)
	}

	// Calculate next cursor (if any)
	var nextCursor *domains.PostGetAllRequest
	if len(postIDs) > pageSize {
		lastPost := posts[pageSize]
		nextCursor = &domains.PostGetAllRequest{
			CursorID:     &lastPost.ID,
			CommentCount: &postResults[pageSize].CommentCount,
		}
		postIDs = postIDs[:pageSize] // Remove the extra record
		posts = posts[:pageSize]     // Remove the extra record
	}
	return posts, nextCursor, err
}
