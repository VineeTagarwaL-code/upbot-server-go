package main

import (
	"log"
	"os"
	"upbot-server-go/database"
	"upbot-server-go/models"
	"upbot-server-go/routes"
	"upbot-server-go/worker"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env")
	}
	database.Connect()
	if err := database.AutoMigrate(&models.User{}, &models.Log{}, &models.Task{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	if err != nil {
		log.Fatal("Error connecting to database")
	}
	r := gin.Default()
	routes.SetupRouter(r)
	PORT := os.Getenv("PORT")
	go worker.NotiWorker()
	go worker.StartPingWorker()

	r.Run(":" + PORT)

}
