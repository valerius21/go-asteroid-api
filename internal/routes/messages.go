package routes

import (
	"github.com/gin-gonic/gin"
)

type Message struct {
	id, uuid, content string
	// createdAt, updatedAt
}

func InitMessages(router *gin.Engine) {

}

func (m Message) create(context *gin.Context) {

}
