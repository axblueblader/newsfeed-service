package service

import (
	"newsfeed-service/domains"
	"newsfeed-service/models"
	"newsfeed-service/storage"
)

type CommentService interface {
	CreateComment(userID string, postID uint, req domains.CreateCommentRequest) (uint, error)
	DeleteComment(userID string, commentID uint) error
}

type commentService struct {
	commentDB storage.CommentDB
}

func NewCommentService(commentDB storage.CommentDB) CommentService {
	return &commentService{commentDB: commentDB}
}

func (s commentService) CreateComment(userID string, postID uint, req domains.CreateCommentRequest) (uint, error) {
	comment, err := s.commentDB.Create(&models.Comment{
		Creator: userID,
		PostID:  postID,
		Content: req.Content,
	})
	if err != nil {
		return 0, err
	}
	return comment.ID, nil
}

func (s commentService) DeleteComment(userID string, commentID uint) error {
	return s.commentDB.Delete(userID, commentID)
}
