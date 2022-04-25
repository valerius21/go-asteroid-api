package note

import (
	"context"
	"github.com/docker/distribution/uuid"
	"github.com/gin-gonic/gin"
	"github.com/pastoapp/astroid-api/internal/orbitdb"
	"log"
	"net/http"
	"time"
)

var DefaultNotes Note

type Note struct {
	ID   uuid.UUID
	Data string // Json Encoded
}

type Transport struct {
	Data string `json:"data" binding:"required"`
}

func init() {
	log.SetPrefix("[middleware/note/note")
}

// TODO: needs merge + refactor
// LET'S GO

func Create(c *gin.Context) {
	var jsonData Transport
	err := c.ShouldBindJSON(&jsonData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := &Note{
		ID:   uuid.Generate(),
		Data: jsonData.Data,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := orbitdb.DefaultDatabase.Put(ctx, map[string]interface{}{
		"_id":  note.ID.String(),
		"data": note.Data,
		"type": "note",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatalf("Could not create note: %v", err)
		return
	}

	_id := resp.GetKey()
	nID, err := uuid.Parse(*_id)

	c.JSON(http.StatusCreated, gin.H{"id": nID.String(), "data": note.Data, "type": "note"})
}

func FindAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	find, err := orbitdb.DefaultDatabase.Query(ctx, func(doc interface{}) (bool, error) {
		return true, nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, find)
}
