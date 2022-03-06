package routes

import (
	"context"
	"log"
	"net/http"

	"berty.tech/go-orbit-db/address"
	"berty.tech/go-orbit-db/iface"
	"github.com/gin-gonic/gin"
	"github.com/pastoapp/astroid-api/internal/orbitdb"
)

var odb = orbitdb.OrbitDB

//const defaultDatabases

type Database struct {
	// id == address
	name      string
	address   address.Address
	storeType string
	store     iface.DocumentStore
}

type databaseDTO struct {
	Name string `json:"name"`
}

func InitDatabases(router *gin.Engine, store iface.Store) {
	log.Println("creating default databases")
	databases := router.Group("/databases")
	messageDB, err := createDatabase("messages")

	if err != nil {
		log.Panicf("%v\n", err)
		// TODO: return error
	}

	userDB, err := createDatabase("users")

	if err != nil {
		log.Panicf("%v\n", err)
		// TODO: return error
	}

	messageDB.createRoute(databases)
	userDB.createRoute(databases)

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
		db, err := createDatabase(d.Name)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": db.toJSON(),
		})
	})
}

// @Summary Create a new database
// @Description Create a new database, it's always a document store
func createDatabase(dbName string) (*Database, error) {
	store, err := orbitdb.CreateStore(dbName)
	if err != nil {
		log.Fatalf("%v\n", err)
		return nil, err
	}

	defer func(store iface.DocumentStore) {
		err := store.Close()
		if err != nil {
			log.Fatalf("Could not close database: %v\n", err)
			return
		}
	}(store)

	return &Database{
		name:      dbName,
		address:   store.Address(),
		storeType: store.Type(),
		store:     store,
	}, nil
}

func (db Database) cacheDB() {
	// TODO: cache db
}

func (db Database) createRoute(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/"+db.name, db.find)
	routerGroup.POST("/"+db.name, func(c *gin.Context) {
		var m MessageDTO
		if err := c.ShouldBindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if m.Text == "" || m.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "message text and ID is required"})
			return
		}

		res, err := db.createItem(m)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})
}

func (db Database) find(c *gin.Context) {
	c.JSON(http.StatusOK, db.toJSON())
}

func (db Database) toJSON() gin.H {
	return gin.H{
		"name":      db.name,
		"address":   db.address.String(),
		"storeType": db.storeType,
	}
}

type MessageDTO struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func (db Database) createItem(m MessageDTO) (gin.H, error) {
	// TODO: add validation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store, err := odb.Docs(ctx, db.address.String(), nil)

	defer func(store iface.DocumentStore) {
		err := store.Close()
		if err != nil {
			log.Fatalf("%v\n", err)
		}
	}(store)

	if err != nil {
		log.Fatalf("%v\n", err)
		return nil, err
	}

	put, err := store.Put(ctx, map[string]interface{}{"_id": m.ID, "text": m.Text})

	return gin.H{
		"key":   put.GetKey(),
		"value": put.GetValue(),
	}, nil
}
