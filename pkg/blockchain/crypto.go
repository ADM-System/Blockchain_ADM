package crypto

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"math/big"
)

// Transaction rappresenta una transazione sulla blockchain
type Transaction struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
	Signature string  `json:"signature"`
	Nonce     int     `json:"nonce"`
	Timestamp int64   `json:"timestamp"`
	PubKey    string  `json:"pubkey"` // Chiave pubblica del mittente in base64
}

// SignTransaction firma una transazione usando la chiave privata fornita
func SignTransaction(tx *Transaction, privKeyBase64 string) error {
	// Decodifica la chiave privata
	privKeyBytes, err := base64.StdEncoding.DecodeString(privKeyBase64)
	if err != nil {
		return fmt.Errorf("errore nella decodifica della chiave privata: %v", err)
	}

	// Converte i bytes in una chiave privata ECDSA
	privKey, err := x509.ParseECPrivateKey(privKeyBytes)
	if err != nil {
		return fmt.Errorf("errore nel parsing della chiave privata: %v", err)
	}

	// Calcola l'hash della transazione
	data := fmt.Sprintf("%s:%s:%f:%d:%d", tx.Sender, tx.Recipient, tx.Amount, tx.Nonce, tx.Timestamp)
	hash := sha256.Sum256([]byte(data))

	// Firma l'hash
	r, s, err := ecdsa.Sign(nil, privKey, hash[:])
	if err != nil {
		return fmt.Errorf("errore nella firma della transazione: %v", err)
	}

	// Componi la firma (r concatenato con s)
	signature := append(r.Bytes(), s.Bytes()...)
	tx.Signature = base64.StdEncoding.EncodeToString(signature)

	// Imposta anche la chiave pubblica nella transazione
	pubBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return fmt.Errorf("errore nella codifica della chiave pubblica: %v", err)
	}
	tx.PubKey = base64.StdEncoding.EncodeToString(pubBytes)

	return nil
}

// VerifySignature verifica che la firma di una transazione sia valida
func VerifySignature(tx Transaction) bool {
	// Se non c'è firma o chiave pubblica, la verifica fallisce
	if tx.Signature == "" || tx.PubKey == "" {
		return false
	}

	// Decodifica la chiave pubblica
	pubkeyBytes, err := base64.StdEncoding.DecodeString(tx.PubKey)
	if err != nil {
		return false
	}

	// Converti i bytes in una chiave pubblica ECDSA
	pubInterface, err := x509.ParsePKIXPublicKey(pubkeyBytes)
	if err != nil {
		return false
	}

	pubkey, ok := pubInterface.(*ecdsa.PublicKey)
	if !ok {
		return false
	}

	// Calcola l'hash della transazione (lo stesso metodo usato per firmare)
	data := fmt.Sprintf("%s:%s:%f:%d:%d", tx.Sender, tx.Recipient, tx.Amount, tx.Nonce, tx.Timestamp)
	hash := sha256.Sum256([]byte(data))

	// Decodifica la firma
	signatureBytes, err := base64.StdEncoding.DecodeString(tx.Signature)
	if err != nil {
		return false
	}

	// Estrai r e s dalla firma (la firma è r concatenato con s)
	if len(signatureBytes) != 64 {
		return false
	}
	r := new(big.Int).SetBytes(signatureBytes[:32])
	s := new(big.Int).SetBytes(signatureBytes[32:])

	// Verifica la firma
	return ecdsa.Verify(pubkey, hash[:], r, s)
}

// GenerateKey genera una nuova coppia di chiavi ECDSA
func GenerateKeys() (string, string, error) {
	// Questa funzione è già implementata in keygen.go
	// Qui si riporta solo la sua firma
	return "", "", nil
}
