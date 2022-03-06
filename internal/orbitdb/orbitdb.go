package orbitdb

import (
	"context"
	"log"
	"net/http"

	odb "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/iface"
	httpapi "github.com/ipfs/go-ipfs-http-client"
)

var OrbitDB iface.OrbitDB

func init() {
	log.Println("Initializing OrbitDB-Context")
	ctx := context.Background()
	// defer cancel()

	// TODO: use a config file
	dbPath := "./data/astroid-api/orbitdb"

	// TODO: make this configurable
	ipfs, err := httpapi.NewURLApiWithClient("localhost:5001", &http.Client{
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
