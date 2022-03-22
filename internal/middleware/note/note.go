package note

import (
	"context"
	"fmt"
	"github.com/docker/distribution/uuid"
	"github.com/pastoapp/astroid-api/internal/orbitdb"
	"log"
)

type Note struct {
	ID   uuid.UUID
	UID  uuid.UUID
	Data interface{}
}

func init() {
	log.SetPrefix("[middleware/note/note] ")
}

func NewNote(uid uuid.UUID, data interface{}) (*Note, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	ctx := context.Background()

	db, err := orbitdb.OpenDatabase(ctx, "notes")
	if err != nil {
		log.Println("Could not open note database")
		return nil, err
	}

	defer func(db *orbitdb.Database) {
		//err := db.Close()
		//if err != nil {
		//	log.Printf("Could not close note database %v\n", nil)
		//}
	}(db)

	//resp, err := db.Create(gin.H{"uid": uid.String(), "data": data}, nil)
	resp, err := db.Create(map[string]interface{}{
		"uid": uid.String(),
		//"data": data,
		"data": "LLorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum Lorem ipsum orem ipsum ",
	}, &orbitdb.DatabaseCreateOptions{
		ID: uuid.Generate().String(),
	})

	if err != nil {
		log.Println("Could not create note")
		return nil, err
	}

	_id := resp["_id"].(string)
	rdata := resp["data"].(map[string]interface{})

	u, err := uuid.Parse(_id)

	if err != nil {
		log.Println("Could not parse note id")
		return nil, err
	}

	return &Note{
		ID:   u,
		UID:  uid,
		Data: rdata["data"],
	}, nil
}

func (_ Note) Find(key string) (*Note, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	ctx := context.Background()

	db, err := orbitdb.OpenDatabase(ctx, "notes")
	if err != nil {
		log.Println("Could not open user database")
		return nil, err
	}

	defer func(db *orbitdb.Database) {
		//err := db.Close()
		//if err != nil {
		//	log.Printf("Could not close note database %v\n", nil)
		//}
	}(db)

	resp, err := db.Read(key)

	if err != nil {
		log.Println("Could not find note")
		return nil, err
	}

	log.Printf("READ NODE %v\n", resp)

	_id := key //resp["_id"].(string)
	rdata := resp["data"].(map[string]interface{})

	u, err := uuid.Parse(_id)

	if err != nil {
		log.Println("Could not parse note id")
		return nil, err
	}

	return &Note{
		ID:   u,
		Data: rdata["data"],
	}, nil
}

func (n Note) FindByUser(uid uuid.UUID) ([]Note, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	ctx := context.Background()

	db, err := orbitdb.OpenDatabase(ctx, "notes")
	if err != nil {
		log.Println("Could not open user database")
		return nil, err
	}

	defer func(db *orbitdb.Database) {
		//err := db.Close()
		//if err != nil {
		//	log.Printf("Could not close note database %v\n", nil)
		//}
	}(db)

	results, err := db.QueryByUID(uid)
	if err != nil {
		return nil, err
	}

	log.Printf("QUERY: %v", results)

	return nil, fmt.Errorf("not implemented")
}
