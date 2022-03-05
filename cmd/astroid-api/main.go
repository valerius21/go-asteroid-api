package main

import (
	"github.com/pastoapp/astroid-api/internal/orbitdb"
	"log"
)

func main() {
	//r := gin.Default()
	//r.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "pong",
	//	})
	//})
	//
	//// init modules
	//routes.InitUsers(r)
	//routes.InitMessages(r)
	//routes.InitAuth(r)

	err := orbitdb.InitOrbitDb()

	if err != nil {
		log.Fatalln(err)
	}

	//r.Run(":3000")

}
