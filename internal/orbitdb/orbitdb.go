package orbitdb

import (
	berty "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/iface"
	"context"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	"log"
	"net/http"
)

var Client berty.OrbitDB

func init() {
	log.SetPrefix("[orbitdb/orbitdb] ")
}

func createUrlHttpApi(ipfsApiURL string) (*httpapi.HttpApi, error) {
	return httpapi.NewURLApiWithClient(ipfsApiURL, &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	})
}

func InitializeOrbitDB(ipfsApiURL, orbitDbDirectory string) (context.CancelFunc, error) {
	// TODO: add config
	// TODO: add other httpapi options
	ctx, cancel := context.WithCancel(context.Background())
	odb, err := NewOrbitDB(ctx, orbitDbDirectory, ipfsApiURL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	Client = odb
	return cancel, nil
}

func NewOrbitDB(ctx context.Context, dbPath, ipfsApiURL string) (iface.OrbitDB, error) {
	coreAPI, err := createUrlHttpApi(ipfsApiURL)

	if err != nil {
		log.Fatalf("Error creating Core API: %v", err)
		return nil, err
	}

	options := &berty.NewOrbitDBOptions{
		Directory: &dbPath,
	}

	return berty.NewOrbitDB(ctx, coreAPI, options)
}
