package orbitdb

import (
	//odb "berty.tech/go-orbit-db"
	"context"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"log"
	"net/http"
)

// TODO: implement orbitdb handlers
func createNode(ctx context.Context, repoPath string) (iface.CoreAPI, error) {
	// open the repo
	repo, err := fsrepo.Open(repoPath)

	if err != nil {
		return nil, err
	}

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption,
		Repo:    repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)

	if err != nil {
		return nil, err
	}

	return coreapi.NewCoreAPI(node)
}

func InitOrbitDb() error {
	log.Println("Creating ODB context")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//_ := "/astroid-api/orbitdb"
	ipfsStorePath := "/home/valerius/.ipfs" //"/astroid-api/ipfs"

	_, err := createNode(ctx, ipfsStorePath)

	if err != nil {
		log.Fatalln("Error creating IPFS Node", err)
		return err
	}

	_, err = httpapi.NewURLApiWithClient("localhost:5001", &http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			DisableKeepAlives: true,
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	//orbit, err := odb.NewOrbitDB(ctx, ipfs, &odb.NewOrbitDBOptions{Directory: &dbPath})
	//
	//if err != nil {
	//	log.Fatalln("Error creating OrbitDB")
	//	return err
	//}
	//
	//identity := orbit.Identity()
	//
	//log.Println(identity.ID)

	return nil
}
