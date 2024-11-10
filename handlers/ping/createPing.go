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

type PingRequest struct {
	Url     string `json:"url" binding:"required,url"`
	WebHook string `json:"webHook"`
}

func CreatePingHandler(c *gin.Context) {
	var pingReq PingRequest

	email, emailExists := c.Get("email")

	if !emailExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid token claims",
			"details": "Email or userId not found in token",
		})
	}
	if err := c.BindJSON(&pingReq); err != nil {
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

	if len(user.Tasks) >= 5 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Task limit reached",
			"details": "You can only have 5 active tasks at a time",
		})
		return
	}
	for _, task := range user.Tasks {
		if task.URL == pingReq.Url {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Task already exists",
				"details": "Task with this URL already exists",
			})
			return
		}
	}
	var webHook *string
	notifyDiscord := false
	if pingReq.WebHook != "" {
		webHook = &pingReq.WebHook
		notifyDiscord = true
	}
	newTask := models.Task{
		URL:           pingReq.Url,
		IsActive:      true,
		WebHook:       webHook,
		NotifyDiscord: notifyDiscord,
		UserID:        user.ID,
	}
	if err := database.DB.Create(&newTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create task",
			"details": err.Error(),
		})

		return
	}

	redisClient := libraries.GetInstance()
	taskMember := fmt.Sprintf("%d|%s", newTask.ID, newTask.URL)
	redisClient.ZAdd(c, "ping_queue", &redis.Z{Score: float64(time.Now().Add(10 * time.Second).Unix()),
		Member: taskMember})
	c.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"url":     pingReq.Url,
		"taskId":  newTask.ID,
	})
}
