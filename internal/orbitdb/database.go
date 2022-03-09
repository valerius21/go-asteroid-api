package orbitdb

import (
	"berty.tech/go-orbit-db/address"
	"berty.tech/go-orbit-db/iface"
	"berty.tech/go-orbit-db/stores/operation"
	"encoding/json"
	"fmt"
	"github.com/docker/distribution/uuid"
	"log"
	"time"
)
import "context"

// Database is the main interface for interacting with OrbitDB
type Database struct {
	Store   *iface.DocumentStore
	Name    string
	Address address.Address
}

type DatabaseCreateOptions struct {
	ID string
}

func init() {
	log.SetPrefix("[orbitdb/database] ")
}

// timeout is used to set the timeout for the database operations
var timeout = 10 * time.Duration(time.Second)

// infinite items to return
var infinite = -1

// OpenDatabase creates or opens a database
func OpenDatabase(ctx context.Context, name string) (*Database, error) {
	if Client == nil {
		log.Fatalf("Client is not initialized")
		return nil, fmt.Errorf("client is not initialized." +
			" Please run orbitdb.InitializeOrbitDB")
	}

	docs, err := Client.Docs(ctx, name, nil)

	if err != nil {
		log.Fatalf("Could not open/create database: %v", err)
		return nil, err
	}

	return &Database{
		Name:    name,
		Store:   &docs,
		Address: docs.Address(),
	}, nil
}

// Create creates a new document in the database
func (d Database) Create(item interface{}, options *DatabaseCreateOptions) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	store := *d.Store
	var put operation.Operation
	var err error

	if options != nil {
		put, err = store.Put(ctx, map[string]interface{}{
			"_id":  options.ID,
			"data": item,
		})
	} else {
		put, err = store.Put(ctx, map[string]interface{}{
			"_id":  uuid.Generate().String(),
			"data": item,
		})
	}

	if err != nil {
		log.Fatalf("Could not create item: %v", err)
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(put.GetValue(), &m)
	if err != nil {
		log.Fatalf("Could not unmarshal item: %v", err)
		return nil, err
	}

	return m, nil
}

// Read reads a document from the database
func (d Database) Read(key string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	store := *d.Store
	err := store.Load(ctx, infinite)

	if err != nil {
		log.Fatalf("Could not load database: %v", err)
		return nil, err
	}

	get, err := store.Get(ctx, key, nil)

	if err != nil {
		log.Fatalf("Could not read item: %v", err)
		return nil, err
	}

	// in case more or less than one item is found
	if len(get) != 1 {
		return make(map[string]interface{}, 0), nil
	}

	item := get[0]

	if err != nil {
		log.Fatalf("Could not unmarshal item: %v", err)
		return nil, err
	}

	return item.(map[string]interface{}), nil
}

// Update updates a document in the database
func (d Database) Update(key string, item interface{}) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	store := *d.Store
	err := store.Load(ctx, infinite)
	if err != nil {
		log.Fatalf("Could not load database: %v", err)
		return nil, err
	}

	// find the item to update
	get, err := store.Get(ctx, key, nil)

	if err != nil {
		log.Fatalf("Error reading item: %v", err)
		return nil, err
	}

	if len(get) != 1 {
		log.Fatalf("Cannot find exactly one item with key %s", key)
		return nil, err
	}

	// update the item
	put, err := store.Put(ctx, map[string]interface{}{
		"_id":  key,
		"data": item,
	})

	if err != nil {
		log.Fatalf("Could not create item: %v", err)
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(put.GetValue(), &m)

	if err != nil {
		log.Fatalf("Could not unmarshal item: %v", err)
		return nil, err
	}

	return m, nil
}

// Delete deletes a document from the database
func (d Database) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	store := *d.Store
	_, err := store.Delete(ctx, key)

	if err != nil {
		log.Fatalf("Could not delete item: %v", err)
		return err
	}

	return nil
}

// Close closes the database
func (d Database) Close() error {
	store := *d.Store
	return store.Close()
}
