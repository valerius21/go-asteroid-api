package routes

import (
	"berty.tech/go-orbit-db/iface"
	"github.com/gin-gonic/gin"
	odb "github.com/pastoapp/astroid-api/internal/orbitdb"
	"log"
	"net/http"
)

// TODO: define default databases

type databaseDTO struct {
	Name string `json:"name"`
}

func InitDatabases(router *gin.Engine, store iface.Store) {
	log.Println("creating default databases")
	databases := router.Group("/databases")
	messageDB, err := odb.CreateDatabase("messages")

	if err != nil {
		log.Panicf("%v\n", err)
		// TODO: return error
	}

	userDB, err := odb.CreateDatabase("users")

	if err != nil {
		log.Panicf("%v\n", err)
		// TODO: return error
	}

	messageDB.CreateRoute(databases)
	userDB.CreateRoute(databases)

	databases.POST("/", func(c *gin.Context) {
		var d databaseDTO
		if err := c.ShouldBindJSON(&d); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if d.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		// TODO: does not create DB
		db, err := odb.CreateDatabase(d.Name)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": db.ToJSON(),
		})
	})
}
