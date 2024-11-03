package main

import (
	"fmt"
	"log"
	"os"
	"upbot-server-go/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env")
	}
	fmt.Print("Hello, World!")
	r := gin.Default()
	routes.SetupRouter(r)
	PORT := os.Getenv("PORT")
	r.Run(":" + PORT)

}
