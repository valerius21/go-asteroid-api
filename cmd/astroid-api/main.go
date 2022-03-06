package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	odb "github.com/pastoapp/astroid-api/internal/orbitdb"
	"github.com/pastoapp/astroid-api/internal/routes"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// init modules
	routes.InitUsers(r)
	routes.InitMessages(r)
	routes.InitAuth(r)

	store, err := odb.CreateStore("test")

	if err != nil {
		panic(err)
	}
	defer store.Close()
	routes.InitDatabases(r, store)

	fmt.Printf("%+v\n", store.Address())

	r.Run(":3000")

}
