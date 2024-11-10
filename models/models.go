package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email string `json:"email" binding:"required,email"`
	Tasks []Task `json:"tasks" gorm:"foreignKey:UserID"`
}
type Task struct {
	gorm.Model
	URL           string  `json:"url" binding:"required"`
	IsActive      bool    `json:"isActive"`
	NotifyDiscord bool    `json:"notifyDiscord"`
	WebHook       *string `json:"webHook" gorm:"default:NULL"`
	UserID        uint    `json:"userId"`
	Logs          []Log   `json:"logs" gorm:"foreignKey:TaskID"`
}
type Log struct {
	gorm.Model
	TaskID      uint      `json:"taskId"`
	Time        time.Time `json:"time"`
	LogResponse string    `json:"logResponse"`
	IsSuccess   bool      `json:"isSuccess"`
}
