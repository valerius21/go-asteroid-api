package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

func Keygen() (string, string) {
	privateK, _ := rsa.GenerateKey(rand.Reader, 4096)

	privkBytes := x509.MarshalPKCS1PrivateKey(privateK)
	privPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkBytes,
		},
	)

	pubkBytes := x509.MarshalPKCS1PublicKey(&privateK.PublicKey)
	pubkPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubkBytes,
	})

	PublicKey := string(pubkPEM)
	PrivateKey := string(privPEM)

	return PublicKey, PrivateKey
}

func SignNonce(privateKey, base64Nonce string) ([]byte, error) {
	nonce, err := base64.StdEncoding.DecodeString(base64Nonce)

	block, _ := pem.Decode([]byte(privateKey))

	privateK, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	sign, err := privateK.Sign(rand.Reader, nonce, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	})

	if err != nil {
		return nil, err
	}

	return sign, nil
}
