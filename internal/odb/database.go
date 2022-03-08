package odb

import (
	"berty.tech/go-orbit-db/address"
	"berty.tech/go-orbit-db/iface"
	"github.com/docker/distribution/uuid"
	"log"
)
import "context"

type Database struct {
	Store   *iface.DocumentStore
	Name    string
	Address address.Address
	Close   func() error // TODO: isolate in function?
}

func init() {
	log.SetPrefix("[database] ")
}

// OpenDatabase creates or opens a database
func OpenDatabase(ctx context.Context, name string) (*Database, error) {
	docs, err := Client.Docs(ctx, name, nil)

	if err != nil {
		log.Fatalf("Could not open/create database: %v", err)
		return nil, err
	}

	return &Database{
		Name:    name,
		Store:   &docs,
		Address: docs.Address(),
		Close:   docs.Close,
	}, nil
}

// TODO: return interfaces instead of []byte

func (d Database) Create(item interface{}) ([]byte, error) {
	// TODO: add timeouts
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := *d.Store
	put, err := store.Put(ctx, map[string]interface{}{
		"_id":  uuid.Generate().String(),
		"data": item,
	})

	if err != nil {
		log.Fatalf("Could not create item: %v", err)
		return nil, err
	}

	return put.GetValue(), nil
}

func (d Database) Read(key string) ([]interface{}, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := *d.Store
	err := store.Load(ctx, -1)

	if err != nil {
		log.Fatalf("Could not load database: %v", err)
		return nil, err
	}

	get, err := store.Get(ctx, key, nil)

	if err != nil {
		log.Fatalf("Could not read item: %v", err)
		return nil, err
	}

	return get, nil
}

func (d Database) Update(key string, item interface{}) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := *d.Store
	err := store.Load(ctx, -1)
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

	return put.GetValue(), nil
}

func (d Database) Delete(key string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	store := *d.Store
	_, err := store.Delete(ctx, key)

	if err != nil {
		log.Fatalf("Could not delete item: %v", err)
		return err
	}

	return nil
}
