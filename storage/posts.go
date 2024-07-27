package storage

import (
	"gorm.io/gorm"
	"newsfeed-service/config"
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
	GetAllWithComments(userID string, cursorID *uint, pageSize int) ([]models.Post, *uint, error)
}

func (dao postDB) Create(post *models.Post) (*models.Post, error) {
	err := dao.db.Create(post).Error
	return post, err
}

func (dao postDB) GetAllWithComments(userID string, cursorID *uint, pageSize int) ([]models.Post, *uint, error) {
	var posts []models.Post
	var postIDs []uint
	type result struct {
		ID           uint
		CommentCount int
	}
	var results []result
	// we can get next cursor or no cursor based on len(result) == pageSize
	limit := pageSize + 1

	args := map[string]interface{}{
		"limit":  limit,
		"userID": userID,
	}

	// without the outer select, it will miss posts with no comments
	// get posts sort by highest comments, then sort by id and created_at for pagination
	// if ID is auto increment, sort by created_at is enough but adding it for clarity in case ID is not auto increment
	// if ID is hashed/random, we need to use created_at combined as cursor
	selectClause := `
            SELECT p.id, COALESCE(c.comment_count, 0) AS comment_count
            FROM posts p
            LEFT JOIN (
                SELECT post_id, COUNT(*) AS comment_count
                FROM comments
                GROUP BY post_id
            ) c ON p.id = c.post_id
          `
	whereClause := `WHERE p.creator = @userID `
	if cursorID != nil {
		whereClause = whereClause + `AND p.id <= @cursorID`
		args["cursorID"] = *cursorID
	}
	orderAndLimit := ` ORDER BY comment_count DESC, p.created_at DESC, p.id DESC LIMIT @limit`

	rawSql := selectClause + whereClause + orderAndLimit

	err := dao.db.Raw(rawSql, args).Scan(&results).Error

	if err != nil {
		return nil, nil, err
	}

	for _, row := range results {
		postIDs = append(postIDs, row.ID)
	}

	// Calculate next cursor (if any)
	var nextCursor *uint
	if len(postIDs) > pageSize {
		nextCursor = &postIDs[pageSize]
		postIDs = postIDs[:pageSize] // Remove the extra record
	} else if len(postIDs) == 0 {
		return posts, nil, nil
	}

	err = dao.db.Preload("Comments", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("created_at desc").Limit(config.Env().CommentsLimit)
	}).
		Where("creator = ? AND id IN ?", userID, postIDs).
		Order("created_at desc").
		Find(&posts).
		Error

	return posts, nextCursor, err
}
