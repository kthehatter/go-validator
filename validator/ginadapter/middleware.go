package ginadapter

import (
	"github.com/gin-gonic/gin"
	"github.com/kthehatter/go-validator/validator"
)

// Middleware creates a Gin middleware for request validation.
func Middleware(options []validator.ValidationOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body gin.H
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"message": "Invalid request body"})
			c.Abort()
			return
		}

		// Run validation and return the first error
		if err := validator.Validate(body, options); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		// Attach the validated body to the context for use in controllers
		c.Set("validatedBody", body)

		// Proceed to the next handler
		c.Next()
	}
}
