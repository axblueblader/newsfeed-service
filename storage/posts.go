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
	GetAllWithComments(userID string) ([]models.Post, error)
}

func (dao postDB) Create(post *models.Post) (*models.Post, error) {
	err := dao.db.Create(post).Error
	return post, err
}

func (dao postDB) GetAllWithComments(userID string) ([]models.Post, error) {
	var posts []models.Post
	err := dao.db.Preload("Comments", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("created_at desc").Limit(config.Env().CommentsLimit)
	}).Where("creator = ?", userID).Find(&posts).Error
	return posts, err
}
