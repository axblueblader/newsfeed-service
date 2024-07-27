package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"newsfeed-service/constants"
	"strings"
)

func BearerTokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if !strings.Contains(token, "Bearer") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid user token"})
			return
		}

		c.Set(constants.UserIdField, strings.TrimPrefix(token, "Bearer "))
		c.Next()
	}
}
