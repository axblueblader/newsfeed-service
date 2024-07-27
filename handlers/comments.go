package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"newsfeed-service/domains"
	"newsfeed-service/service"
	"strconv"
)

type CommentsHandler struct {
	CommentService service.CommentService
}

func (h CommentsHandler) CreateComment(c *gin.Context) {
	var req domains.CreateCommentRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := GetUserID(c)
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	commentID, err := h.CommentService.CreateComment(userID, uint(postID), req)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{"ID": commentID})
}

func (h CommentsHandler) DeleteComment(c *gin.Context) {
	userID := GetUserID(c)
	commentID, err := strconv.Atoi(c.Param("commentID"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.CommentService.DeleteComment(userID, uint(commentID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ID": commentID})
}
