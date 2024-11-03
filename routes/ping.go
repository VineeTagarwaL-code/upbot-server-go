package routes

import (
	"upbot-server-go/handlers/ping"

	"github.com/gin-gonic/gin"
)

func pingRouter(r *gin.RouterGroup) {
	r.POST("/create", ping.CreatePingHandler)
}
