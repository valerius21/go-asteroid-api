package routes

import (
	"github.com/gin-gonic/gin"
	"log"
)

type User struct {
	_id       string
	publicKey string
	nonce     string
	isAdmin   bool
	// messages, createdAt, updatedAt
}

func InitUsers(router *gin.Engine) {
	log.Println("Test users")
	users := router.Group("/users")
	users.GET("/", User{}.findAll)
}

func (u User) find(context *gin.Context) {

}

func (u User) findAll(context *gin.Context) {
	context.JSON(200, gin.H{
		"user": "dummy",
	})
}

func (u User) create(context *gin.Context) {

}

func (u User) update(context *gin.Context) {

}

func (u User) requestToken(context *gin.Context) {

}
