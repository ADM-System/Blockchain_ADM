package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Struttura per rappresentare una transazione
type Transaction struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

// Struttura per rappresentare un blocco
type Block struct {
	Index        int           // Posizione del blocco nella catena
	Timestamp    string        // Quando è stato creato il blocco
	Transactions []Transaction // Dati memorizzati nel blocco
	PrevHash     string        // Hash del blocco precedente
	Hash         string        // Hash del blocco corrente
	Nonce        int           // Numero usato per il mining
}

// Struttura per rappresentare la blockchain
type Blockchain struct {
	Chain        []Block // La catena di blocchi
	Difficulty   int     // Difficoltà del mining
	MiningReward float64 // Ricompensa per il mining
}

// Aggiungi una variabile globale per le transazioni in sospeso
var pendingTransactions []Transaction

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

// Funzione per il mining di un blocco
func (bc *Blockchain) MineBlock(transactions []Transaction, miner string) {
	prevBlock := bc.Chain[len(bc.Chain)-1]

	// Aggiungi la transazione di ricompensa per il miner
	rewardTx := Transaction{
		Sender:    "SYSTEM",
		Recipient: miner,
		Amount:    bc.MiningReward,
	}
	transactions = append(transactions, rewardTx)

	newBlock := NewBlock(prevBlock.Index+1, transactions, prevBlock.Hash)

	// Definisci la difficoltà
	difficulty := 4                                                 // Numero di zeri all'inizio dell'hash
	fmt.Printf("Inizio mining del blocco #%d...\n", newBlock.Index) // Messaggio di inizio mining
	for {
		hash := calculateHash(newBlock)
		if strings.HasPrefix(hash, strings.Repeat("0", difficulty)) {
			newBlock.Hash = hash
			bc.Chain = append(bc.Chain, newBlock)
			fmt.Printf("Blocco #%d minato con successo! Hash: %s\n", newBlock.Index, newBlock.Hash) // Messaggio di successo
			break
		}
		newBlock.Nonce++ // Incrementa il nonce e riprova

		// Aggiungi un ritardo di 1 millisecondi
		time.Sleep(1 / 2)
	}
}

// Funzione per aggiungere un nuovo blocco alla blockchain
func (bc *Blockchain) AddBlock(transactions []Transaction) {
	bc.MineBlock(transactions, "") // Usa la funzione di mining
}

// Funzione per verificare la validità della blockchain
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		prevBlock := bc.Chain[i-1]

		// Controlla se l'hash del blocco corrente è corretto
		if currentBlock.Hash != calculateHash(currentBlock) {
			return false
		}

		// Controlla se il prevHash del blocco corrente è corretto
		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
	}
	return true
}

// Funzione per salvare la blockchain su disco
func (bc *Blockchain) SaveToFile(filename string) error {
	data, err := json.Marshal(bc)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// Funzione per caricare la blockchain da un file
func LoadBlockchainFromFile(filename string) (Blockchain, error) {
	var bc Blockchain
	data, err := os.ReadFile(filename)
	if err != nil {
		return bc, err
	}
	err = json.Unmarshal(data, &bc)
	return bc, err
}

// Funzione per calcolare il saldo di un indirizzo
func (bc *Blockchain) GetBalance(address string) float64 {
	balance := 0.0

	// Scorri tutti i blocchi
	for _, block := range bc.Chain {
		// Scorri tutte le transazioni nel blocco
		for _, tx := range block.Transactions {
			// Se l'indirizzo è il destinatario, aggiungi l'importo
			if tx.Recipient == address {
				balance += tx.Amount
			}
			// Se l'indirizzo è il mittente, sottrai l'importo
			if tx.Sender == address {
				balance -= tx.Amount
			}
		}
	}

	return balance
}

// Funzione principale
func main() {
	blockchain, err := LoadBlockchainFromFile("blockchain.json")
	if err != nil {
		blockchain = NewBlockchain(4) // Inizializza la blockchain se non esiste
	}

	// Inizializza il mempool vuoto
	pendingTransactions = []Transaction{}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// Endpoint per aggiungere una transazione al mempool
	http.HandleFunc("/addTransaction", func(w http.ResponseWriter, r *http.Request) {
		var tx Transaction
		if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
			http.Error(w, "Transazione non valida", http.StatusBadRequest)
			return
		}
		pendingTransactions = append(pendingTransactions, tx)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Transazione aggiunta al mempool."))
	})

	// Endpoint per minare tutte le transazioni in sospeso
	http.HandleFunc("/mine", func(w http.ResponseWriter, r *http.Request) {
		if len(pendingTransactions) == 0 {
			http.Error(w, "Nessuna transazione da minare.", http.StatusBadRequest)
			return
		}
		// Ricevi il nome del miner dal body (JSON)
		var data struct {
			Miner string `json:"miner"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil || data.Miner == "" {
			http.Error(w, "Nome del miner mancante o non valido.", http.StatusBadRequest)
			return
		}
		blockchain.MineBlock(pendingTransactions, data.Miner)
		pendingTransactions = []Transaction{} // Svuota il mempool dopo il mining
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Blocco minato con successo!"))
	})

	http.HandleFunc("/blockchain", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(blockchain.Chain)
	})

	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		isValid := blockchain.IsValid()
		if isValid {
			w.Write([]byte("La blockchain è valida."))
		} else {
			w.Write([]byte("La blockchain non è valida."))
		}
	})

	http.HandleFunc("/mempool", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(pendingTransactions)
	})

	http.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Indirizzo non specificato", http.StatusBadRequest)
			return
		}
		balance := blockchain.GetBalance(address)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"address": address,
			"balance": balance,
		})
	})

	defer func() {
		if err := blockchain.SaveToFile("blockchain.json"); err != nil {
			fmt.Println("Errore nel salvataggio della blockchain:", err)
		} else {
			fmt.Println("Blockchain salvata con successo.")
		}
	}()

	fmt.Println("Server in esecuzione su http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
