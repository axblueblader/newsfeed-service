package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"newsfeed-service/domains"
	"newsfeed-service/services"
)

type PostsHandler struct {
	PostService services.PostService
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
	// default page size of 10
	req := domains.PostGetAllRequest{
		PageSize: 10,
	}
	err := c.BindQuery(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	posts, err := h.PostService.GetPostsWithComments(GetUserID(c), req.CursorID, req.PageSize)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}
