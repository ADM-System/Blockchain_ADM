package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Struttura per rappresentare una transazione
type Transaction struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

// Block rappresenta un singolo blocco nella blockchain
type Block struct {
	Index        int           // Posizione del blocco nella catena
	Timestamp    string        // Quando è stato creato il blocco
	Transactions []Transaction // Dati memorizzati nel blocco
	PrevHash     string        // Hash del blocco precedente
	Hash         string        // Hash del blocco corrente
	Nonce        int           // Numero usato per il mining
}

// Blockchain è una lista di blocchi
type Blockchain struct {
	Chain        []Block
	Difficulty   int // Difficoltà del mining
	MiningReward int // Ricompensa per il mining
}

// Funzione per creare un nuovo blocco
func NewBlock(index int, transactions []Transaction, prevHash string) Block {
	block := Block{
		Index:        index,
		Timestamp:    time.Now().String(),
		Transactions: transactions,
		PrevHash:     prevHash,
		Hash:         "",
		Nonce:        0,
	}
	return block
}

// Funzione per calcolare l'hash di un blocco
func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%s%s%d", block.Index, block.Timestamp, block.Transactions, block.PrevHash, block.Nonce)
	hash := sha256.New()
	hash.Write([]byte(record))
	hashInBytes := hash.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}

// Funzione per inizializzare la blockchain
func NewBlockchain(difficulty int) Blockchain {
	genesisBlock := NewBlock(0, []Transaction{}, "")
	genesisBlock.Hash = calculateHash(genesisBlock)
	return Blockchain{
		Chain:        []Block{genesisBlock},
		Difficulty:   difficulty,
		MiningReward: 100, // Ricompensa di default
	}
}

// Funzione per aggiungere un nuovo blocco alla blockchain
func (bc *Blockchain) AddBlock(transactions []Transaction) {
	prevBlock := bc.Chain[len(bc.Chain)-1]
	newBlock := NewBlock(prevBlock.Index+1, transactions, prevBlock.Hash)
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
			Transactions []Transaction `json:"transactions"`
		}
		json.NewDecoder(r.Body).Decode(&newBlock)
		blockchain.AddBlock(newBlock.Transactions)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/blockchain", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(blockchain.Chain)
	})

	fmt.Println("Server in esecuzione su http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
