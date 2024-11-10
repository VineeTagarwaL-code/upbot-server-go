package ping

import (
	"fmt"
	"net/http"
	"time"
	"upbot-server-go/database"
	"upbot-server-go/libraries"
	"upbot-server-go/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type ReactivatePingRequest struct {
	TaskId int `json:"taskId" binding:"required"`
}

func ReactivatePingHandler(c *gin.Context) {
	var pingReq ReactivatePingRequest

	email, emailExists := c.Get("email")
	if !emailExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid token claims",
			"details": "Email not found in token",
		})
		return
	}

	if err := c.BindJSON(&pingReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	var user models.User
	err := database.DB.Preload("Tasks").First(&user, "email = ?", email).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "User not found",
			"details": err.Error(),
		})
		return
	}

	var task models.Task
	err = database.DB.First(&task, "id = ? AND user_id = ?", pingReq.TaskId, user.ID).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Task not found for this user",
			"details": err.Error(),
		})
		return
	}

	task.IsActive = true
	if err := database.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to reactivate task",
			"details": err.Error(),
		})
		return
	}

	redisClient := libraries.GetInstance()
	taskMember := fmt.Sprintf("%d|%s", task.ID, task.URL)
	nextPingTime := time.Now().Add(10 * time.Second).Unix()
	_, err = redisClient.ZAdd(c, "ping_queue", &redis.Z{
		Score:  float64(nextPingTime),
		Member: taskMember,
	}).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to add task to ping queue",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task reactivated successfully",
		"taskId":  task.ID,
		"url":     task.URL,
	})
}
