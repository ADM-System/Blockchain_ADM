package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Block rappresenta un singolo blocco nella blockchain
type Block struct {
	Index     int    // Posizione del blocco nella catena
	Timestamp string // Quando è stato creato il blocco
	Data      string // Dati memorizzati nel blocco
	PrevHash  string // Hash del blocco precedente
	Hash      string // Hash del blocco corrente
	Nonce     int    // Numero usato per il mining
}

// Blockchain è una lista di blocchi
type Blockchain struct {
	Chain        []Block
	Difficulty   int // Difficoltà del mining
	MiningReward int // Ricompensa per il mining
}

// Funzione per creare un nuovo blocco
func NewBlock(index int, data, prevHash string) Block {
	block := Block{
		Index:     index,
		Timestamp: time.Now().String(),
		Data:      data,
		PrevHash:  prevHash,
		Hash:      "",
		Nonce:     0,
	}
	return block
}

// Funzione per calcolare l'hash di un blocco
func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%s%s%d", block.Index, block.Timestamp, block.Data, block.PrevHash, block.Nonce)
	hash := sha256.New()
	hash.Write([]byte(record))
	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}

// Funzione per verificare se l'hash soddisfa la difficoltà
func (block *Block) hasValidHash(difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(block.Hash, prefix)
}

// Funzione per il mining di un blocco
func (block *Block) mineBlock(difficulty int) {
	for !block.hasValidHash(difficulty) {
		block.Nonce++
		block.Hash = calculateHash(*block)
	}
}

// Funzione per inizializzare la blockchain
func NewBlockchain(difficulty int) Blockchain {
	genesisBlock := NewBlock(0, "Genesis Block", "")
	genesisBlock.mineBlock(difficulty)
	return Blockchain{
		Chain:        []Block{genesisBlock},
		Difficulty:   difficulty,
		MiningReward: 100, // Ricompensa di default
	}
}

// Funzione per verificare se un blocco è valido
func (block *Block) IsValid() bool {
	// Verifica che l'hash calcolato corrisponda all'hash memorizzato
	calculatedHash := calculateHash(*block)
	return calculatedHash == block.Hash
}

// Funzione per verificare l'intera blockchain
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		previousBlock := bc.Chain[i-1]

		// Verifica che l'hash del blocco corrente sia valido
		if !currentBlock.IsValid() {
			return false
		}

		// Verifica che il blocco corrente punti correttamente al blocco precedente
		if currentBlock.PrevHash != previousBlock.Hash {
			return false
		}
	}
	return true
}

// Funzione per aggiungere un nuovo blocco alla blockchain
func (bc *Blockchain) AddBlock(data string) error {
	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := NewBlock(prevBlock.Index+1, data, prevBlock.Hash)

	// Mining del nuovo blocco
	newBlock.mineBlock(bc.Difficulty)

	// Verifica che il nuovo blocco sia valido
	if !newBlock.IsValid() {
		return fmt.Errorf("il nuovo blocco non è valido")
	}

	bc.Chain = append(bc.Chain, newBlock)
	return nil
}

func main() {
	// Inizializziamo una nuova blockchain con difficoltà 4
	blockchain := NewBlockchain(4)

	// Aggiungiamo alcuni blocchi
	err := blockchain.AddBlock("First block after Genesis")
	if err != nil {
		fmt.Printf("Errore nell'aggiunta del blocco: %v\n", err)
		return
	}

	err = blockchain.AddBlock("Second block after Genesis")
	if err != nil {
		fmt.Printf("Errore nell'aggiunta del blocco: %v\n", err)
		return
	}

	// Verifichiamo l'integrità della blockchain
	if blockchain.IsValid() {
		fmt.Println("La blockchain è valida!")
	} else {
		fmt.Println("La blockchain non è valida!")
	}

	// Mostriamo i blocchi della blockchain
	for _, block := range blockchain.Chain {
		fmt.Printf("Block #%d\n", block.Index)
		fmt.Printf("Timestamp: %s\n", block.Timestamp)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("PrevHash: %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Nonce: %d\n\n", block.Nonce)
	}
}
