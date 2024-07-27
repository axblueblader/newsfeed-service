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
