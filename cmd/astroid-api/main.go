package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

var (
	ipfsURL       = "http://localhost:5001"
	orbitDbDir    = "./data/orbitdb"
	defaultStores = []string{"users", "notes"}
)

func main() {
	// main database context
	//ctx, cancel := context.WithCancel(context.Background())
	//defer  cancel()
	//// create a new orbitdb instance
	//cancelODB, err := odb.InitializeOrbitDB(ipfsURL, orbitDbDir)
	//defer cancelODB() // cancel the orbitdb context

	//for _, store := range defaultStores {
	//	database, err := odb.OpenDatabase(ctx, store)
	//	if err != nil {
	//        log.Fatal(err)
	//	}
	//	database.
	//}

	// gin server
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	err := r.Run(":3000")

	if err != nil {
		log.Panicf("Error starting server: %v", err)
	}

}
