package note

import (
	"context"
	"github.com/gin-gonic/gin"
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

	item := gin.H{
		"Hi": "Mom",
	}

	t.Run("should create a note", func(t *testing.T) {
		note, err := NewNote(item)

		if err != nil {
			t.Fatalf("could not create note %v\n", err)
		}

		if note.Data == nil {
			t.Fatalf("note item is empty")
		}

		if note.Data.(map[string]interface{})["Hi"] != "Mom" {
			t.Fatalf("note item is not correct")
		}
	})

	t.Run("should create and find some notes", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			note, err := NewNote(item)

			if err != nil {
				t.Fatalf("could not create note %v\n", err)
			}

			if note.Data == nil {
				t.Fatalf("note item is empty")
			}

			if note.Data.(map[string]interface{})["Hi"] != "Mom" {
				t.Fatalf("note item is not correct")
			}

			qNote, err := Note{}.Find(note.ID.String())

			if err != nil {
				t.Fatalf("could not find note %v\n", err)
			}

			if qNote.Data.(map[string]interface{})["Hi"] != "Mom" {
				t.Fatalf("note item is not correct")
			}

			if qNote.ID != note.ID {
				t.Fatalf("note id is not correct")
			}
		}
	})

}
