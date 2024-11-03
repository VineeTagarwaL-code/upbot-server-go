package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandleFuncWithError func(c *gin.Context) error

func ErrorWrapper(handler HandleFuncWithError) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(c); err != nil {
			fmt.Print("error occured")
			log.Printf("Error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort() // Stop further processing of the request
		}
	}
}
