package routes

import (
	"github.com/gin-gonic/gin"
	"testing"
)
import "gopkg.in/h2non/gock.v1"

// TODO: maybe take a look into mockery package
func TestUserRoutes(t *testing.T) {
	// Disable gock
	defer gock.Off()
	// gin server
	r := gin.Default()

	t.Run("should test the /ping endpoint", func(t *testing.T) {
		r.GET("/ping", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"message": "pong",
			})
		})
		err := r.Run(":8080")
		if err != nil {
			t.Fatalf("Error while running the server: %v", err)
		}
	})

	t.Run("should create a user", func(t *testing.T) {
		t.Skip("Not implemented")
	})
}
