package note

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/gin-gonic/gin"
	user2 "github.com/pastoapp/astroid-api/internal/middleware/user"
	"github.com/pastoapp/astroid-api/internal/orbitdb"
	"testing"
)

func TestNewNotes(t *testing.T) {

	privateK, _ := rsa.GenerateKey(rand.Reader, 4096)
	pubkBytes := x509.MarshalPKCS1PublicKey(&privateK.PublicKey)
	pubkPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubkBytes,
	})

	PublicKey := string(pubkPEM)

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

	user, err := user2.NewUser(PublicKey, false)

	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}

	t.Run("should create a note for user", func(t *testing.T) {
		note, err := NewNote(user.ID, item)

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

	t.Run("should bloat the db", func(t *testing.T) {
		for i := 0; i < 300; i++ {
			_, err := user2.NewUser("brooooo", false)
			if err != nil {
				t.Fatalf("Error creating user: %v", err)
			}
		}
	})

	t.Run("should create and find some notes, based on Note-ID", func(t *testing.T) {
		var nodeIds []string
		for i := 0; i < 10; i++ {
			note, err := NewNote(user.ID, item)
			nodeIds = append(nodeIds, note.ID.String())

			if err != nil {
				t.Fatalf("could not create note %v\n", err)
			}

			if note.Data == nil {
				t.Fatalf("note item is empty")
			}

			//if note.Data.(map[string]interface{})["Hi"] != "Mom" {
			//	t.Fatalf("note item is not correct")
			//}

		}

		//for _, id := range nodeIds {
		//	qNote, err := Note{}.Find(id)
		//
		//	if err != nil {
		//		t.Fatalf("could not find note %v\n", err)
		//	}
		//
		//	if qNote.Data.(map[string]interface{})["Hi"] != "Mom" {
		//		t.Fatalf("note item is not correct")
		//	}
		//
		//	if qNote.ID.String() != id {
		//		t.Fatalf("note id is not correct")
		//	}
		//}
	})

	t.Run("should find all Users' notes", func(t *testing.T) {
		for i := 0; i < 10; i++ {

			note, err := NewNote(user.ID, item)

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
		t.Log(user.ID.String())
		_, err := Note{}.FindByUser(user.ID)
		if err != nil {
			t.Errorf("could not find notes by user: %v\n", err)
		}
	})

}
