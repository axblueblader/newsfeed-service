package handlers

import (
	"github.com/gin-gonic/gin"
	"newsfeed-service/constants"
)

func GetUserID(c *gin.Context) string {
	return c.GetString(constants.UserIdField)
}
