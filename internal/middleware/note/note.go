package note

import (
	"context"
	"github.com/docker/distribution/uuid"
	"github.com/pastoapp/astroid-api/internal/orbitdb"
	"log"
	"time"
)

type Note struct {
	ID   uuid.UUID
	Data interface{}
}

func init() {
	log.SetPrefix("[middleware/note/note")
}

func NewNote(data interface{}) (*Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := orbitdb.OpenDatabase(ctx, "notes")
	if err != nil {
		log.Fatalln("Could not open user database")
		return nil, err
	}

	defer func(db *orbitdb.Database) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Could not close note database %v\n", nil)
		}
	}(db)

	resp, err := db.Create(data, nil)

	if err != nil {
		log.Fatalln("Could not create note")
		return nil, err
	}

	_id := resp["_id"].(string)
	rdata := resp["data"].(map[string]interface{})

	u, err := uuid.Parse(_id)

	if err != nil {
		log.Fatalln("Could not parse note id")
		return nil, err
	}

	return &Note{
		ID:   u,
		Data: rdata,
	}, nil
}

func (n Note) Find(key string) (*Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := orbitdb.OpenDatabase(ctx, "notes")
	if err != nil {
		log.Fatalln("Could not open user database")
		return nil, err
	}

	defer func(db *orbitdb.Database) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Could not close note database %v\n", nil)
		}
	}(db)

	resp, err := db.Read(key)

	if err != nil {
		log.Fatalln("Could not create note")
		return nil, err
	}

	_id := resp["_id"].(string)
	rdata := resp["Data"].(map[string]interface{})

	u, err := uuid.Parse(_id)

	if err != nil {
		log.Fatalln("Could not parse note id")
		return nil, err
	}

	return &Note{
		ID:   u,
		Data: rdata,
	}, nil
}
