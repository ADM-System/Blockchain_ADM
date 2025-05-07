package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
}

func main() {
	// Inserisci qui i dati della transazione
	tx := Transaction{
		Sender:    "admin",
		Recipient: "bob",
		Amount:    10,
	}

	// Inserisci qui la chiave privata (base64)
	privStr := "INCOLLA_LA_TUA_CHIAVE_PRIVATA_BASE64"
	privBytes, _ := base64.StdEncoding.DecodeString(privStr)
	priv, _ := x509.ParseECPrivateKey(privBytes)

	// Firma la transazione
	data := fmt.Sprintf("%s:%s:%f", tx.Sender, tx.Recipient, tx.Amount)
	hash := sha256.Sum256([]byte(data))
	r, s, _ := ecdsa.Sign(nil, priv, hash[:])
	signature := append(r.Bytes(), s.Bytes()...)
	signatureStr := base64.StdEncoding.EncodeToString(signature)

	// Stampa la firma da incollare nel campo "signature"
	fmt.Println("Firma (da incollare nel campo 'signature'):\n", signatureStr)
}
