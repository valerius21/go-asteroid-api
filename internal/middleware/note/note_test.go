package note

import (
	"context"
	"github.com/pastoapp/astroid-api/internal/orbitdb"
	"testing"
)

func TestNewNotes(t *testing.T) {
	cancelFunc, err := orbitdb.InitializeOrbitDB("http://localhost:5001", t.TempDir())
	if err != nil {
		t.Fatalf("Error initializing OrbitDB: %v", err)
	}
	defer cancelFunc()
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Skip("Skipping test")

}
