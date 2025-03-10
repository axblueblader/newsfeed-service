package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Server is healthy",
	})
}
