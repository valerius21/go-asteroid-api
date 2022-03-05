package main

import (
	"fmt"

	odb "github.com/pastoapp/astroid-api/internal/orbitdb"
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

	//r.Run(":3000")
	store, err := odb.CreateStore("test")
	if err != nil {
		panic(err)
	}
	defer store.Close()

	fmt.Printf("%+v\n", store)

}
