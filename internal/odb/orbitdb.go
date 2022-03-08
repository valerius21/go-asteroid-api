package odb

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
	log.SetPrefix("[orbitdb] ")
}

func createCoreAPI(ipfsApiURL string) (*httpapi.HttpApi, error) {
	return httpapi.NewURLApiWithClient(ipfsApiURL, &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	})
}

func NewOrbitDB(ctx context.Context, dbPath string) (iface.OrbitDB, error) {
	// TODO: add config
	ipfsApiURL := "http://localhost:5001"
	coreAPI, err := createCoreAPI(ipfsApiURL)

	if err != nil {
		log.Fatalf("Error creating Core API: %v", err)
		return nil, err
	}

	options := &berty.NewOrbitDBOptions{
		Directory: &dbPath,
	}

	return berty.NewOrbitDB(ctx, coreAPI, options)
}
