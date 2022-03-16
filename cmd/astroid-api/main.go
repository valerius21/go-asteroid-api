package main

import (
	"context"
	"github.com/pastoapp/astroid-api/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
	odb "github.com/pastoapp/astroid-api/internal/orbitdb"
)

var (
	ipfsURL       = "http://localhost:5001"
	orbitDbDir    = "./data/orbitdb"
	defaultStores = []string{"users", "notes"}
)

func main() {
	// main database context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// create a new orbitdb instance
	cancelODB, err := odb.InitializeOrbitDB(ipfsURL, orbitDbDir)
	defer cancelODB() // cancel the orbitdb context

	userDB, err := odb.OpenDatabase(ctx, defaultStores[0])
	if err != nil {
		log.Fatal(err)
	}

	_, err = odb.OpenDatabase(ctx, defaultStores[1])
	if err != nil {
		log.Fatal(err)
	}

	// gin server
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	routes.InitUsers(r, userDB)

	err = r.Run(":3000")

	if err != nil {
		log.Panicf("Error starting server: %v", err)
	}
}
