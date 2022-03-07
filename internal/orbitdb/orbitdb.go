package orbitdb

import (
	"berty.tech/go-orbit-db/address"
	"context"
	"fmt"
	"github.com/docker/distribution/uuid"
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

// CreateStore creates a new document store
func CreateStore(name string) (iface.DocumentStore, error) {
	store, err := OrbitDB.Docs(context.Background(), name, nil)

	return store, err
}

// CreateDatabase creates a new database
//  Create a new database, it's always a document store
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

// CreateRoute creates a new route for the database relative to the given gin.RouterGroup
func (db Database) CreateRoute(routerGroup *gin.RouterGroup) {
	routerGroup.GET("/"+db.name, db.findAll)
	routerGroup.POST("/"+db.name, func(c *gin.Context) {
		var m MessageDTO
		if err := c.ShouldBindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Data Field Malformat. "+
				"Please check your request parameters, %v\n", err.Error())})
			return
		}
		m.ID = uuid.Generate().String()
		if m.Data == nil {
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

// Find returns all items in the database
func (db Database) findAll(c *gin.Context) {
	var fn OperationFn = func(ctx context.Context, store iface.DocumentStore) (gin.H, error){
		filter := func(doc interface{}) (bool, error) {
			return true, nil
		}
		res, err := store.Query(ctx, filter, nil)
	}
	c.JSON(http.StatusOK, db.ToJSON())
}

// ToJSON returns the database as a JSON object
func (db Database) ToJSON() gin.H {
	return gin.H{
		"name":      db.name,
		"address":   db.address.String(),
		"storeType": db.storeType,
	}
}

// MessageDTO is a message data transport object
type MessageDTO struct {
	ID   string
	Data map[string]interface{} `json:"data"`
}

type OperationFn func(ctx context.Context, store iface.DocumentStore) (gin.H, error)

// OpenAndDo performs an operation on the database and returns the result
// This function is used to open the database and close it after the operation
// is done.
// @param operationFn is the operation to perform on the database.
// @returns the result of the operation
func (db Database) OpenAndDo(fn OperationFn, options *iface.CreateDBOptions) (gin.H, error) {
	// TODO: add validation
	// TODO: add example
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store, err := OrbitDB.Docs(ctx, db.address.String(), options)

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
	// operation
	res, err := fn(ctx, store)
	if err != nil {
		log.Fatalf("%v\n", err)
		return nil, err
	}

	return res, nil
}

// CreateItem creates a new item in the database
func (db Database) createItem(m MessageDTO) (gin.H, error) {
	fn := func(ctx context.Context, store iface.DocumentStore) (gin.H, error) {
		put, err := store.Put(ctx, map[string]interface{}{"_id": m.ID, "data": m.Data})
		if err != nil {
			return nil, err
		}
		return gin.H{
			"key":   put.GetKey(),
			"value": put.GetValue(),
		}, nil
	}
	return db.OpenAndDo(fn, nil)
}
