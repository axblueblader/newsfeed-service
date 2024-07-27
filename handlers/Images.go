package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"newsfeed-service/config"
	"newsfeed-service/domains"
	"newsfeed-service/storage"
)

type ImagesHandler struct {
	ObjectStorage storage.ObjectStorage
}

func (h ImagesHandler) GenerateSignedUrl(c *gin.Context) {
	userId := c.GetString("userId")
	imageUuid := uuid.New().String()
	c.JSON(http.StatusOK, domains.SignedUrlResponse{
		SignedUrl: h.ObjectStorage.GenerateSignedUrl(config.Env().PostImageBucketName, userId+"/"+imageUuid),
	})
}
func (h ImagesHandler) ProcessPostImageUploaded(c *gin.Context) {
	var req domains.ImageUploadedRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Do image conversion based
	// convertImage(bucket,path)

	// Replace the image in the current url
	// replaceImage(bucket,path)

	// Now the post will always serve the processed image

	c.JSON(http.StatusOK, req)
}
