package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	pingGroup := r.Group("/ping")

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server is healthy"})
	})
	pingRouter(pingGroup)

}
