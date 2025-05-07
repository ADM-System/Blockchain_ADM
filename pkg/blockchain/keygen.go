package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

func main() {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		panic(err)
	}
	pubStr := base64.StdEncoding.EncodeToString(pubBytes)
	fmt.Println("Chiave pubblica (da usare nel campo 'pubkey'):\n", pubStr)
	// La chiave privata va tenuta segreta! Per demo, la stampiamo in base64
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		panic(err)
	}
	privStr := base64.StdEncoding.EncodeToString(privBytes)
	fmt.Println("Chiave privata (usala solo per firmare, NON condividerla!):\n", privStr)
}
