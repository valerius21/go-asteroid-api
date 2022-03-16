package routes

import (
	"github.com/gin-gonic/gin"
)

// TODO: define default databases

type databaseDTO struct {
	Name string `json:"name"`
}

func InitDatabases(router *gin.Engine) {

}
