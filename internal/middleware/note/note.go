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

// Note entity, holding a note in the OrbitDB.
type Note struct {
	ID   uuid.UUID
	Data string // Json Encoded
}

// Transport holds the provided information in a JSON input
type Transport struct {
	// Data is the parsed json input
	Data string `json:"data" binding:"required"`
}

// init runs at module initialization.
func init() {
	log.SetPrefix("[middleware/note/note")
}

// TODO: needs merge + refactor
// LET'S GO

// Create takes a *gin.Context argument, which binds the associated route to create a Note in the OrbitDB instance.
func Create(c *gin.Context) {

	var jsonData Transport
	err := c.ShouldBindJSON(&jsonData)

	// reject malformed input.
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// initialize the new Note.
	note := &Note{
		ID:   uuid.Generate(),
		Data: jsonData.Data,
	}

	// Attempt to add a Note to the OrbitDB, giving it a timeout of 10 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := orbitdb.DefaultDatabase.Put(ctx, map[string]interface{}{
		"_id":  note.ID.String(),
		"data": note.Data,
		"type": "note",
	})

	// if query was not successful, throw an error with a response object.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatalf("Could not create note: %v", err)
		return
	}

	// Otherwise, format the OrbitDB response for using it in the HTTP response, so that the requester gets notified.
	_id := resp.GetKey()
	nID, err := uuid.Parse(*_id)

	c.JSON(http.StatusCreated, gin.H{"id": nID.String(), "data": note.Data, "type": "note"})
}

// FindAll returns all Note Objects from the OrbitDB. It exists for the purpose of PoC. Requesting single Note objects
// needs a dedicated filter function.
func FindAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Query the ODB with a filter function; currently allowing every object to be returned (PoC).
	find, err := orbitdb.DefaultDatabase.Query(ctx, func(doc interface{}) (bool, error) {
		return true, nil
	})

	// if the query fails, response with an error.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return found Note objects
	c.JSON(http.StatusOK, find)
}
