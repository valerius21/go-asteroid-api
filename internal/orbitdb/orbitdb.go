package orbitdb

import (
	//odb "berty.tech/go-orbit-db"
	"context"
	"fmt"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"io/ioutil"
	"log"
)

// TODO: implement orbitdb handlers

func createRepo(repoPath string) error {
	cfg, err := config.Init(ioutil.Discard, 2048)

	if err != nil {
		return err
	}

	cfg.Experimental.FilestoreEnabled = true
	cfg.Experimental.UrlstoreEnabled = true
	cfg.Experimental.Libp2pStreamMounting = true
	cfg.Experimental.P2pHttpProxy = true
	cfg.Experimental.StrategicProviding = true

	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return fmt.Errorf("failed to init node %s", err)
	}

	return nil
}

func createNode(ctx context.Context, repoPath string) (iface.CoreAPI, error) {
	// create the repo
	err := createRepo(repoPath)

	if err != nil {
		return nil, err
	}
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
	ipfsStorePath := "$HOME/.ipfs" //"/astroid-api/ipfs"

	_, err := createNode(ctx, ipfsStorePath)

	if err != nil {
		log.Fatalln("Error creating IPFS Node", err)
		return err
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
