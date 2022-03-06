package orbitdb

import (
	"berty.tech/go-orbit-db/address"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"

	odb "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/iface"
	ipfsHttpApi "github.com/ipfs/go-ipfs-http-client"
)

var OrbitDB iface.OrbitDB

type Database struct {
	// id == address
	name      string
	address   address.Address
	storeType string
	store     iface.DocumentStore
}

func init() {
	log.Println("Initializing OrbitDB-Context")
	ctx := context.Background()
	// defer cancel()

	// TODO: use a config file
	dbPath := "./data/asteroid-api/orbitdb"

	// TODO: make this configurable
	ipfs, err := ipfsHttpApi.NewURLApiWithClient("localhost:5001", &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	})

	if err != nil {
		log.Panicf("Error creating IPFS client: %s", err)
	}

	orbit, err := odb.NewOrbitDB(ctx, ipfs, &odb.NewOrbitDBOptions{Directory: &dbPath})

	if err != nil {
		log.Panicf("Error creating OrbitDB: %s", err)
	}

	identity := orbit.Identity()

	log.Printf("Initialized OrbitDB with ID: %s", identity.ID)
	OrbitDB = orbit
}

func CreateStore(name string) (iface.DocumentStore, error) {
	store, err := OrbitDB.Docs(context.Background(), name, nil)

	return store, err
}

// CreateDatabase @Summary Create a new database
// @Description Create a new database, it's always a document store
func CreateDatabase(dbName string) (*Database, error) {
	store, err := CreateStore(dbName)
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

func (db Database) CreateRoute(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/"+db.name, db.find)
	routerGroup.POST("/"+db.name, func(c *gin.Context) {
		var m MessageDTO
		if err := c.ShouldBindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Data Field Malformat. "+
				"Please check your request parameters, %v\n", err.Error())})
			return
		}
		if m.Data == nil || m.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "message data and ID are required"})
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
	c.JSON(http.StatusOK, db.ToJSON())
}

func (db Database) ToJSON() gin.H {
	return gin.H{
		"name":      db.name,
		"address":   db.address.String(),
		"storeType": db.storeType,
	}
}

type MessageDTO struct {
	ID   string                 `json:"id"`
	Data map[string]interface{} `json:"data"`
}

func (db Database) createItem(m MessageDTO) (gin.H, error) {
	// TODO: add validation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store, err := OrbitDB.Docs(ctx, db.address.String(), nil)

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

	put, err := store.Put(ctx, map[string]interface{}{"_id": m.ID, "data": m.Data})

	if err != nil {
		log.Fatalf("%v\n", err)
		return nil, err
	}

	return gin.H{
		"key":   put.GetKey(),
		"value": put.GetValue(),
	}, nil
}
