package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Block rappresenta un singolo blocco nella blockchain
type Block struct {
	Index     int    // Posizione del blocco nella catena
	Timestamp string // Quando è stato creato il blocco
	Data      string // Dati memorizzati nel blocco
	PrevHash  string // Hash del blocco precedente
	Hash      string // Hash del blocco corrente
}

// Blockchain è una lista di blocchi
type Blockchain struct {
	Chain []Block
}

// Funzione per creare un nuovo blocco
func NewBlock(index int, data, prevHash string) Block {
	block := Block{
		Index:     index,
		Timestamp: time.Now().String(),
		Data:      data,
		PrevHash:  prevHash,
		Hash:      "",
	}
	// Calcoliamo l'hash del blocco
	block.Hash = calculateHash(block)
	return block
}

// Funzione per calcolare l'hash di un blocco
func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%s%s", block.Index, block.Timestamp, block.Data, block.PrevHash)
	hash := sha256.New()
	hash.Write([]byte(record))
	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}

// Funzione per inizializzare la blockchain
func NewBlockchain() Blockchain {
	// Il primo blocco è chiamato "Genesis Block"
	genesisBlock := NewBlock(0, "Genesis Block", "")
	return Blockchain{Chain: []Block{genesisBlock}}
}

// Funzione per aggiungere un nuovo blocco alla blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := NewBlock(prevBlock.Index+1, data, prevBlock.Hash)
	bc.Chain = append(bc.Chain, newBlock)
}

func main() {
	// Inizializziamo una nuova blockchain
	blockchain := NewBlockchain()

	// Aggiungiamo alcuni blocchi
	blockchain.AddBlock("First block after Genesis")
	blockchain.AddBlock("Second block after Genesis")

	// Mostriamo i blocchi della blockchain
	for _, block := range blockchain.Chain {
		fmt.Printf("Block #%d\n", block.Index)
		fmt.Printf("Timestamp: %s\n", block.Timestamp)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("PrevHash: %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n\n", block.Hash)
	}
}
