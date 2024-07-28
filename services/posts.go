package services

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
	GetPostsWithComments(userID string, req domains.PostGetAllRequest) (*domains.PostsPagedResult, error)
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

func (s postService) GetPostsWithComments(userID string, req domains.PostGetAllRequest) (*domains.PostsPagedResult, error) {
	posts, nextCursor, err := s.postDB.GetAllWithComments(userID, req)
	if err != nil {
		return nil, err
	}
	postsWithComments := make([]domains.PostWithComments, 0)
	for _, post := range posts {
		comments := make([]domains.Comment, 0)
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
	result := &domains.PostsPagedResult{
		Posts:    postsWithComments,
		PageSize: req.PageSize,
	}
	if nextCursor != nil {
		result.NextCursor = nextCursor.CursorID
		result.CommentCount = nextCursor.CommentCount
	}
	return result, nil
}
