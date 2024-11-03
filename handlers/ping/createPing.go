package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingRequest struct {
	Url      string `json:"url"`
	Interval int16  `json:"interval"`
}

func CreatePingHandler(c *gin.Context) {

	var pingReq PingRequest

	err := c.BindJSON(&pingReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ping created successfully", "url": pingReq.Url})
}
