package storage

import (
	"gorm.io/gorm"
	"newsfeed-service/models"
)

type commentDB struct {
	db *gorm.DB
}

func NewCommentDB(db *gorm.DB) CommentDB {
	return &commentDB{db: db}
}

type CommentDB interface {
	Create(comment *models.Comment) (*models.Comment, error)
	Delete(userID string, commentID uint) error
}

func (dao commentDB) Create(comment *models.Comment) (*models.Comment, error) {
	err := dao.db.Create(comment).Error
	return comment, err
}

func (dao commentDB) Delete(userID string, commentID uint) error {
	return dao.db.Where("creator = ? AND id = ?", userID, commentID).Delete(&models.Comment{}).Error
}
