package ping

import (
	"net/http"
	"upbot-server-go/database"
	"upbot-server-go/models"

	"github.com/gin-gonic/gin"
)

func GetPingsHandler(c *gin.Context) {

	email, emailExists := c.Get("email")

	if !emailExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid token claims",
			"details": "Email or userId not found in token",
		})
	}
	user := models.User{}
	if err := database.DB.Preload("Tasks.Logs").Find(&user, "email = ?", email).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pings fetched successfully",
		"additional": gin.H{
			"pings": user.Tasks,
		},
	})
}
