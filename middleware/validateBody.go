package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateBody takes a function that returns a new instance of the object to validate.
func ValidateBody(objFactory func() interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new instance of the object for each request
		obj := objFactory()

		// Bind JSON to the newly created object
		if err := c.BindJSON(obj); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			c.Abort()
			return
		}

		// Validate the object
		if err := validate.Struct(obj); err != nil {
			errors := make(map[string]string)
			for _, err := range err.(validator.ValidationErrors) {
				errors[err.Field()] = err.Tag()
			}
			c.JSON(http.StatusBadRequest, gin.H{"validationErrors": errors})
			c.Abort()
			return
		}
		c.Set("validatedBody", obj) // Optionally set the validated object in context
		c.Next()
	}
}
