package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"newsfeed-service/domains"
	"newsfeed-service/service"
)

type PostsHandler struct {
	PostService service.PostService
}

func (h PostsHandler) CreatePost(c *gin.Context) {
	var req domains.PostCreateRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	postID, err := h.PostService.CreatePost(GetUserID(c), req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ID": postID})
}

func (h PostsHandler) RetrievePostWithComments(c *gin.Context) {
	posts, err := h.PostService.GetPostsWithComments(GetUserID(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}
