package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
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

// Funzione per inizializzare la blockchain
func NewBlockchain(difficulty int) Blockchain {
	genesisBlock := NewBlock(0, "Genesis Block", "")
	genesisBlock.Hash = calculateHash(genesisBlock)
	return Blockchain{
		Chain:        []Block{genesisBlock},
		Difficulty:   difficulty,
		MiningReward: 100, // Ricompensa di default
	}
}

// Funzione per aggiungere un nuovo blocco alla blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := NewBlock(prevBlock.Index+1, data, prevBlock.Hash)
	newBlock.Hash = calculateHash(newBlock)
	bc.Chain = append(bc.Chain, newBlock)
}

// Funzione principale
func main() {
	blockchain := NewBlockchain(4) // Inizializza la blockchain

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/addBlock", func(w http.ResponseWriter, r *http.Request) {
		var newBlock struct {
			Data string `json:"data"`
		}
		json.NewDecoder(r.Body).Decode(&newBlock)
		blockchain.AddBlock(newBlock.Data)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/blockchain", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(blockchain.Chain)
	})

	fmt.Println("Server in esecuzione su http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
