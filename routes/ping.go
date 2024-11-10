package routes

import (
	"upbot-server-go/handlers/ping"
	"upbot-server-go/middleware"

	"github.com/gin-gonic/gin"
)

type PingRequest struct {
	Url string `json:"url" validate:"required,url"`
}

func pingRouter(r *gin.RouterGroup) {
	r.POST("/create", middleware.TokenValidator, ping.CreatePingHandler)
	r.DELETE("/delete", middleware.TokenValidator, ping.DeletePingHandler)
	r.GET("/getall", middleware.TokenValidator, ping.GetPingsHandler)
	r.PATCH("/reactivate", middleware.TokenValidator, ping.ReactivatePingHandler)
}
