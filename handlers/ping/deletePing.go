package ping

import (
	"net/http"
	"upbot-server-go/database"
	"upbot-server-go/models"

	"github.com/gin-gonic/gin"
)

type DeletePingRequest struct {
	TaskId uint `json:"taskId" binding:"required"`
}

func DeletePingHandler(c *gin.Context) {
	var delPingReq DeletePingRequest

	email, emailExists := c.Get("email")

	if !emailExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid token claims",
			"details": "Email or userId not found in token",
		})
	}
	if err := c.BindJSON(&delPingReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	var user models.User
	err := database.DB.Preload("Tasks").Find(&user, "email = ?", email).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}

	for _, task := range user.Tasks {
		if task.ID == delPingReq.TaskId {
			if err := database.DB.Delete(&task).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to delete task",
					"details": err.Error(),
				})
				return
			}

			if err := database.DB.Where("task_id = ?", task.ID).Delete(&models.Log{}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to delete logs",
					"details": err.Error(),
				})
				return
			}
			if err := database.DB.Model(&user).Association("Tasks").Delete(&task); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to delete user relation",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Task deleted successfully",
				"taskId":  delPingReq.TaskId,
			})
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "Task not found",
		"details": "Task with this URL not found",
	})

}
