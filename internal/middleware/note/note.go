package note

import (
	"github.com/docker/distribution/uuid"
	"log"
)

type Note struct {
	ID   uuid.UUID
	Data interface{}
}

func init() {
	log.SetPrefix("[middleware/note/note")
}

// TODO: needs merge + refactor
