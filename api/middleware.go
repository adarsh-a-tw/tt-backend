package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func adminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Api-Key") != "" && c.GetHeader("X-Api-Key") != os.Getenv("ADMIN_TOKEN") {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
