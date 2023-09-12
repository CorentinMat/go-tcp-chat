package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

func (c *client) encrypt(msg string) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// The public key is a part of the *rsa.PrivateKey struct
	publicKey := privateKey.PublicKey

	fmt.Println("private / public : ", privateKey, publicKey)
}
