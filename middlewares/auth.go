package middlewares

import "github.com/gin-gonic/gin"

const dummy_token = "secret"

func Auth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token != dummy_token {
		c.AbortWithStatusJSON(403, "Invalid token")
	}
}
