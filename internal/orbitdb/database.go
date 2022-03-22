package orbitdb

import (
	"berty.tech/go-orbit-db/address"
	"berty.tech/go-orbit-db/iface"
	"berty.tech/go-orbit-db/stores/operation"
	"encoding/json"
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

// Create creates a new document in the database
func Create(item interface{}, options *DatabaseCreateOptions) (map[string]interface{}, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var put operation.Operation
	var err error

	if options != nil {
		put, err = DefaultDatabase.Put(ctx, map[string]interface{}{
			"_id":  options.ID,
			"data": item,
		})
	} else {
		put, err = DefaultDatabase.Put(ctx, map[string]interface{}{
			"_id":  uuid.Generate().String(),
			"data": item,
		})
	}

	if err != nil {
		log.Fatalf("Could not create item 1: %v", err)
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
func Read(key string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	get, err := DefaultDatabase.Get(ctx, key, nil)

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
func Update(key string, item interface{}) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// find the item to update
	get, err := DefaultDatabase.Get(ctx, key, nil)

	if err != nil {
		log.Fatalf("Error reading item: %v", err)
		return nil, err
	}

	if len(get) != 1 {
		log.Fatalf("Cannot find exactly one item with key %s", key)
		return nil, err
	}

	// update the item
	put, err := DefaultDatabase.Put(ctx, map[string]interface{}{
		"_id":  key,
		"data": item,
	})

	if err != nil {
		log.Fatalf("Could not create item 2: %v", err)
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
func Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err := DefaultDatabase.Delete(ctx, key)
	return err
}
