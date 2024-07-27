package service

import (
	"newsfeed-service/domains"
	"newsfeed-service/models"
	"newsfeed-service/storage"
)

type postService struct {
	postDB storage.PostDB
}

func NewPostService(postDB storage.PostDB) PostService {
	return &postService{postDB: postDB}
}

type PostService interface {
	CreatePost(userID string, req domains.PostCreateRequest) (uint, error)
	GetPostsWithComments(userID string) ([]domains.PostWithComments, error)
}

func (s postService) CreatePost(userID string, req domains.PostCreateRequest) (uint, error) {
	post, err := s.postDB.Create(&models.Post{
		Caption:  req.Caption,
		Creator:  userID,
		ImageUrl: req.ImageUrl,
	})
	if err != nil {
		return 0, err
	}
	return post.ID, nil
}

func (s postService) GetPostsWithComments(userID string) ([]domains.PostWithComments, error) {
	posts, err := s.postDB.GetAllWithComments(userID)
	var postsWithComments []domains.PostWithComments
	for _, post := range posts {
		var comments []domains.Comment
		for _, comment := range post.Comments {
			comments = append(comments, domains.Comment{
				ID:        comment.ID,
				Content:   comment.Content,
				Creator:   comment.Creator,
				CreatedAt: comment.CreatedAt,
			})
		}

		postsWithComments = append(postsWithComments, domains.PostWithComments{
			ID:        post.ID,
			Caption:   post.Caption,
			ImageUrl:  post.ImageUrl,
			Creator:   post.Creator,
			CreatedAt: post.CreatedAt,
			Comments:  comments,
		})
	}
	return postsWithComments, err
}
