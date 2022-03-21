package user

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/docker/distribution/uuid"
	"github.com/pastoapp/astroid-api/internal/orbitdb"
	"log"
	"time"
)

type User struct {
	ID        uuid.UUID
	PublicKey string
	Nonce     string
	IsAdmin   bool
	CreatedAt int64
	UpdatedAt int64
	//note TODO: add note
}

func init() {
	log.SetPrefix("[middleware/user/user] ")
}

func NewUser(publicKey string, isAdmin bool) (*User, error) {
	nonce, err := GenerateNonce()
	if err != nil {
		log.Fatalln("Failed to generate Nonce")
		return nil, err
	}

	user := &User{
		ID:        uuid.Generate(),
		PublicKey: publicKey,
		// TODO: REGENERATE NONCE EVERY TIME AN AUTH SUCCESSFULLY HAPPENS
		// base64 encoded nonce
		Nonce:     nonce,
		IsAdmin:   isAdmin,
		CreatedAt: time.Now().UTC().Unix(),
		UpdatedAt: time.Now().UTC().Unix(),
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

	resp, err := db.Create(*user, nil)

	if err != nil {
		log.Fatalln("Could not create user")
		return nil, err
	}

	_id := resp["_id"].(string)

	newID, err := uuid.Parse(_id)

	if err != nil {
		log.Fatalln("Could not parse UUID")
		return nil, err
	}

	return &User{
		ID:        newID,
		PublicKey: user.PublicKey,
		Nonce:     user.Nonce,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func GenerateNonce() (string, error) {
	key := [64]byte{}
	_, err := rand.Read(key[:])
	if err != nil {
		log.Fatalln("Failed to generate random key")
		return "", err
	}

	msgHash := sha256.New()
	_, err = msgHash.Write(key[:])
	if err != nil {
		log.Fatalln("Failed to hash key")
		return "", err
	}
	return base64.StdEncoding.EncodeToString(msgHash.Sum(nil)), nil
}

func (u User) Login() (string, error) {
	// TODO: implement
	// TODO: return JWT
	return "", fmt.Errorf("not implemented")
}

// RefreshNonce updates the user nonce
func (u User) RefreshNonce() error {
	nonce, err := GenerateNonce()
	if err != nil {
		log.Fatalln("Failed to generate Nonce")
		return err
	}
	u.Nonce = nonce
	return nil
}

func (u User) VerifyUser(signature string) error {

	block, _ := pem.Decode([]byte(u.PublicKey))
	if block == nil {
		return fmt.Errorf("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)

	if err != nil {
		return fmt.Errorf("failed to parse DER encoded public key: %s\n", err.Error())
	}

	nonce, err := base64.StdEncoding.DecodeString(u.Nonce)

	if err != nil {
		return fmt.Errorf("failed to decode nonce: %s\n", err.Error())
	}

	return rsa.VerifyPSS(pub, crypto.SHA256, nonce, []byte(signature), nil)
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

	find, err := db.Read(key)

	if err != nil {
		log.Fatalln("Could not find user")
		return nil, err
	}

	id, err := uuid.Parse(find["_id"].(string))

	if err != nil {
		log.Fatalln("Could not parse user id")
		return nil, err
	}

	data := find["data"].(map[string]interface{})

	ca := data["CreatedAt"].(float64)
	ua := data["UpdatedAt"].(float64)

	return &User{
		ID:        id,
		PublicKey: data["PublicKey"].(string),
		Nonce:     data["Nonce"].(string),
		IsAdmin:   data["IsAdmin"].(bool),
		CreatedAt: int64(ca),
		UpdatedAt: int64(ua),
	}, nil
}
