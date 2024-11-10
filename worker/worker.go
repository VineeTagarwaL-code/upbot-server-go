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
			time.Sleep(1 * time.Minute)
			continue
		}
		for _, task := range tasks {

			parts := strings.SplitN(task, "|", 2)
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
			resp, err := http.Get(url)
			if err != nil {
				nextPing := time.Now().Add(10 * time.Minute).Unix()
				_, err = redisClient.ZAdd(context.Background(), "ping_queue", &redis.Z{
					Score:  float64(nextPing),
					Member: url,
				}).Result()
				newLog := models.Log{
					LogResponse: "Failed to ping URL",
					Time:        time.Now(),
					TaskID:      uint(taskId),
					IsSuccess:   false,
				}
				if err := database.DB.Create(&newLog).Error; err != nil {
					log.Printf("Error creating log: %v", err)
				}
				if err != nil {
					log.Printf("Error rescheduling URL %s: %v", url, err)
				}

				continue
			}
			resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				nextPing := time.Now().Add(10 * time.Minute).Unix()
				_, err = redisClient.ZAdd(context.Background(), "ping_queue", &redis.Z{
					Score:  float64(nextPing),
					Member: url,
				}).Result()
				if err != nil {
					log.Printf("Error rescheduling URL %s: %v", url, err)
				}
				newLog := models.Log{
					LogResponse: "Successfully pinged - 200",
					Time:        time.Now(),
					TaskID:      uint(taskId),
					IsSuccess:   true,
				}
				if err := database.DB.Create(&newLog).Error; err != nil {
					log.Printf("Error creating log: %v", err)
				}

			} else {
				nextPing := time.Now().Add(10 * time.Minute).Unix()
				_, err = redisClient.ZAdd(context.Background(), "ping_queue", &redis.Z{
					Score:  float64(nextPing),
					Member: url,
				}).Result()
				newLog := models.Log{
					LogResponse: "Ping Unsuccesfull - " + strconv.Itoa(resp.StatusCode),
					Time:        time.Now(),
					TaskID:      uint(taskId),
					IsSuccess:   false,
				}
				if err := database.DB.Create(&newLog).Error; err != nil {
					log.Printf("Error creating log: %v", err)
				}
				if err != nil {
					log.Printf("Error rescheduling URL %s: %v", url, err)
				}
			}
		}

		time.Sleep(10 * time.Minute)
	}
}
