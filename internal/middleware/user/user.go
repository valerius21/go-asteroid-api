package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/docker/distribution/uuid"
	"github.com/pastoapp/astroid-api/internal/orbitdb"
	"log"
	"time"
)

type User struct {
	ID        uuid.UUID
	publicKey string
	nonce     string
	isAdmin   bool
	createdAt int64
	updatedAt int64
	//notes TODO: add notes
}

func init() {
	log.SetPrefix("[middleware/user/user] ")
}

func GenerateNonce() (string, error) {
	key := [64]byte{}
	_, err := rand.Read(key[:])
	if err != nil {
		log.Fatalln("Failed to generate random key")
		return "", err
	}

	return base64.URLEncoding.EncodeToString([]byte(fmt.Sprint(key))), nil
}

func NewUser(publicKey string, isAdmin bool) (*User, error) {

	nonce, err := GenerateNonce()
	if err != nil {
		log.Fatalln("Failed to generate nonce")
		return nil, err
	}

	user := &User{
		ID:        uuid.Generate(),
		publicKey: publicKey,
		nonce:     nonce,
		isAdmin:   isAdmin,
		createdAt: time.Now().UTC().Unix(),
		updatedAt: time.Now().UTC().Unix(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := orbitdb.OpenDatabase(ctx, "users")

	if err != nil {
		log.Fatalln("Could not open user database")
		return nil, err
	}

	defer func(db *orbitdb.Database) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Could not close user database %v\n", user)
		}
	}(db)

	_, err = db.Create(user, &orbitdb.DatabaseCreateOptions{ID: user.ID.String()})

	if err != nil {
		log.Fatalln("Could not create user")
		return nil, err
	}

	return user, nil
}

func (u User) Login() (string, error) {
	// TODO: implement
	// TODO: return JWT
	return "", fmt.Errorf("not implemented")
}

func Find(key string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := orbitdb.OpenDatabase(ctx, "users")

	if err != nil {
		log.Fatalln("Could not open user database")
		return nil, err
	}

	defer func(db *orbitdb.Database) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Could not close user database %v\n", key)
		}
	}(db)

	_, err = db.Read(key)

	if err != nil {
		log.Fatalln("Could not find user")
		return nil, err
	}
	return nil, fmt.Errorf("not implemented")
}
