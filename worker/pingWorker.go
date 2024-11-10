package worker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"upbot-server-go/database"
	"upbot-server-go/libraries"
	"upbot-server-go/models"

	"github.com/go-redis/redis/v8"
)

func StartPingWorker() {
	for {
		now := time.Now().Unix()
		redisClient := libraries.GetInstance()
		tasks, err := redisClient.ZRangeByScore(context.Background(), "ping_queue", &redis.ZRangeBy{
			Min: "-inf",
			Max: fmt.Sprintf("%d", now),
		}).Result()
		if err != nil {
			log.Printf("Error fetching from queue: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		for _, task := range tasks {

			parts := strings.SplitN(task, "|", 2)
			fmt.Println(parts)
			fmt.Println(len(parts))
			if len(parts) != 2 {
				log.Printf("Invalid task format: %s", task)
				continue
			}
			taskIdStr, url := parts[0], parts[1]
			taskId, err := strconv.Atoi(taskIdStr)
			if err != nil {
				log.Printf("Invalid task ID: %s", taskIdStr)
				continue
			}
			timeNow := time.Now()
			resp, err := http.Get(url)
			timeSince := time.Since(timeNow).Milliseconds()
			taskMember := fmt.Sprintf("%d|%s", taskId, url)

			if err != nil || resp.StatusCode != http.StatusOK {
				newLog := models.Log{
					LogResponse: "Failed to ping URL",
					Time:        time.Now(),
					TimeTake:    int64(timeSince),
					TaskID:      uint(taskId),
					IsSuccess:   false,
				}
				if err := database.DB.Create(&newLog).Error; err != nil {
					log.Printf("Error creating log: %v", err)
				}

				var task models.Task
				if err := database.DB.First(&task, taskId).Error; err == nil {
					task.FailCount++
					if task.FailCount >= 2 {
						task.IsActive = false
						database.DB.Save(&task)
						redisClient.ZRem(context.Background(), "ping_queue", taskMember)
						log.Printf("Task %d has failed more than 2 times. Marked as inactive and removed from queue.", taskId)
					} else {
						database.DB.Model(&task).Update("fail_count", task.FailCount)
						nextPing := time.Now().Add(10 * time.Second).Unix()
						_, err = redisClient.ZAdd(context.Background(), "ping_queue", &redis.Z{
							Score:  float64(nextPing),
							Member: taskMember,
						}).Result()
						if err != nil {
							log.Printf("Error rescheduling URL %s: %v", url, err)
						}
					}
				} else {
					log.Printf("Error fetching task: %v", err)
				}
				continue
			}
			if resp.StatusCode == http.StatusOK {
				taskMember := fmt.Sprintf("%d|%s", taskId, url)
				nextPing := time.Now().Add(10 * time.Second).Unix()
				_, err = redisClient.ZAdd(context.Background(), "ping_queue", &redis.Z{
					Score:  float64(nextPing),
					Member: taskMember,
				}).Result()
				if err != nil {
					log.Printf("Error rescheduling URL %s: %v", url, err)
				}
				fmt.Println(timeSince)
				newLog := models.Log{
					LogResponse: "Successfully pinged - 200",
					Time:        time.Now(),
					TimeTake:    int64(timeSince),
					TaskID:      uint(taskId),
					IsSuccess:   true,
				}
				if err := database.DB.Create(&newLog).Error; err != nil {
					log.Printf("Error creating log: %v", err)
				}
			}
			resp.Body.Close()
		}

		time.Sleep(10 * time.Second)
	}
}
