package orbitdb

import (
	"encoding/json"
	"testing"
)
import "context"

func closeDb(db *Database, t *testing.T) {
	err := db.Close()
	if err != nil {
		t.Fatalf("Error closing database: %v", err)
	}
}

func TestNewDatabase(t *testing.T) {
	cancelFunc, err := InitializeOrbitDB("http://localhost:5001", t.TempDir())
	if err != nil {
		t.Fatalf("Error initializing OrbitDB: %v", err)
	}
	defer cancelFunc()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("should create a new database", func(t *testing.T) {
		db, err := OpenDatabase(ctx, "create-test")
		defer closeDb(db, t)

		if db == nil {
			t.Errorf("expected database to be created")
		}

		if err != nil {
			t.Errorf("error creating database: %s", err)
		}

		if db.Name != "create-test" {
			t.Errorf("expected database name to be 'testdb', got %s", db.Name)
		}

		if db.Store == nil {
			t.Errorf("expected database store to be created")
		}

		if len(db.Address.String()) < 10 {
			t.Errorf("expected database address to be set")
		}
	})

	item := map[string]interface{}{"Hi": "mom"}

	// TODO: test validations
	t.Run("should create a new item in the database", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		db, err := OpenDatabase(ctx, "rw-global-test")
		defer closeDb(db, t)

		if err != nil {
			t.Errorf("error creating database: %s", err)
		}
		resp, err := db.Create(item)

		if err != nil {
			t.Errorf("error adding item: %s", err)
		}

		if len(resp) == 0 {
			t.Errorf("expected response to have length > 0")
		}

		if resp == nil {
			t.Errorf("expected response to be returned")
		}
	})

	t.Run("should update an item in the database", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		db, err := OpenDatabase(ctx, "rw-global-test")

		if err != nil {
			t.Errorf("error creating database: %s", err)
		}

		// prepare the item
		resp, err := db.Create(item)

		if err != nil {
			t.Errorf("error adding item: %s", err)
		}

		m := make(map[string]interface{})
		err = json.Unmarshal(resp, &m)

		if err != nil {
			t.Errorf("error unmarshalling response: %s", err)
		}

		_id := m["_id"].(string)
		data := m["data"].(map[string]interface{})
		value := data["Hi"].(string)

		if value != "mom" {
			t.Errorf("expected value to be 'mom', got %s", value)
		}

		if _id == "" {
			t.Errorf("expected id to be set")
		}

		closeDb(db, t)
		cancel()

		// update the item
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
		db, err = OpenDatabase(ctx, "rw-global-test")
		defer closeDb(db, t)

		if err != nil {
			t.Errorf("error creating database: %s", err)
		}

		updated, err := db.Update(_id, map[string]interface{}{"Hi": "dad"})

		if err != nil {
			t.Errorf("error reading item: %s", err)
		}

		if updated == nil {
			t.Errorf("expected item to be returned")
		}
		m = make(map[string]interface{})
		err = json.Unmarshal(updated, &m)

		if m["_id"].(string) != _id {
			t.Errorf("expected id to be %s, got %s", _id, m["_id"].(string))
		}

		if m["data"].(map[string]interface{})["Hi"] != "dad" {
			t.Errorf("expected value to be 'dad', got %s", m["data"].(map[string]interface{})["Hi"])
		}
	})

	t.Run("should delete an item with the specified key", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		db, err := OpenDatabase(ctx, "rw-global-test")
		defer closeDb(db, t)

		if err != nil {
			t.Errorf("error creating database: %s", err)
		}

		resp, err := db.Create(item)

		m := make(map[string]interface{})
		err = json.Unmarshal(resp, &m)
		_id := m["_id"].(string)

		err = db.Delete(_id)

		if err != nil {
			t.Errorf("error deleting item: %s", err)
		}

		get, err := db.Read(_id)

		if len(get) != 0 {
			t.Errorf("expected item to be deleted")
		}
	})
}
