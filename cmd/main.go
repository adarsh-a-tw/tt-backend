package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var port = 8080

func main() {
	addr := fmt.Sprintf(":%d", port)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(addr)
}
