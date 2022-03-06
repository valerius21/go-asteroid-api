package routes

import (
	"context"
	"fmt"
	"log"

	"berty.tech/go-orbit-db/address"
	"berty.tech/go-orbit-db/iface"
	"github.com/gin-gonic/gin"
	"github.com/pastoapp/astroid-api/internal/orbitdb"
)

var odb = orbitdb.OrbitDB

type Database struct {
	// __id      string
	name      string
	address   address.Address
	storeType string
}

func InitDatabases(router *gin.Engine, store iface.Store) {
	log.Println("Test databases")
	databases := router.Group("/databases")
	databases.GET("/", Database{
		name:      store.DBName(),
		address:   store.Address(),
		storeType: store.Type(),
	}.findAll)
	// databases.POST("/", Database{}.create)
}

func (db Database) findAll(c *gin.Context) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	store, err := odb.Open(ctx, db.address.String(), nil)

	if err != nil {
		log.Fatalf("Error opening store: %s", err)
	}

	c.JSON(200, gin.H{
		"dbData": fmt.Sprintf("%+v", store),
	})
}

func (db Database) create(c *gin.Context) {

}
