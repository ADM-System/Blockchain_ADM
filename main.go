package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

var peers = []string{
	"http://localhost:8081", // aggiungi qui altri peer se vuoi
}

// Struttura per rappresentare una transazione
type Transaction struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
	Signature string  `json:"signature"`
	Nonce     int     `json:"nonce"`
	Timestamp int64   `json:"timestamp"` // Unix timestamp
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

var lastMiningTime = make(map[string]time.Time)

const miningCooldown = 10 * time.Second // tempo minimo tra due mining dello stesso miner
const mempoolTxMaxAge = 3600            // secondi (1 ora)

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
	// Transazione speciale per dare 100 coins ad "admin"
	genesisTx := Transaction{
		Sender:    "SYSTEM",
		Recipient: "admin",
		Amount:    100,
	}
	genesisBlock := NewBlock(0, []Transaction{genesisTx}, "")
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

func syncWithPeers() {
	for {
		for _, peer := range peers {
			resp, err := http.Get(peer + "/blockchain")
			if err != nil {
				continue
			}
			defer resp.Body.Close()
			var peerChain []Block
			if err := json.NewDecoder(resp.Body).Decode(&peerChain); err != nil {
				continue
			}
			if len(peerChain) > len(blockchain.Chain) {
				// Recupera le transazioni orfane PRIMA di sostituire la catena
				orphans := recoverOrphanTransactions(blockchain.Chain, peerChain, pendingTransactions)
				blockchain.Chain = peerChain
				// Rimetti le orfane nel mempool (evita duplicati)
				for _, tx := range orphans {
					alreadyInMempool := false
					for _, pending := range pendingTransactions {
						if pending.Signature == tx.Signature && pending.Nonce == tx.Nonce {
							alreadyInMempool = true
							break
						}
					}
					if !alreadyInMempool {
						pendingTransactions = append(pendingTransactions, tx)
					}
				}
				fmt.Println("Catena aggiornata da peer:", peer, "e transazioni orfane recuperate:", len(orphans))
			}
		}
		time.Sleep(10 * time.Second)
	}
}

// Funzione per pulire il mempool dalle transazioni troppo vecchie
func cleanMempool() {
	now := time.Now().Unix()
	newPending := []Transaction{}
	for _, tx := range pendingTransactions {
		if now-tx.Timestamp <= mempoolTxMaxAge {
			newPending = append(newPending, tx)
		}
	}
	pendingTransactions = newPending
}

// recoverOrphanTransactions identifica le transazioni che erano nella vecchia catena
// ma non nella nuova, e le restituisce per reinserirle nel mempool
func recoverOrphanTransactions(oldChain []Block, newChain []Block, currentMempool []Transaction) []Transaction {
	// Mappa per tenere traccia delle transazioni nella nuova catena
	txInNewChain := make(map[string]bool)

	// Mappa le transazioni nella vecchia catena
	txInOldChain := make(map[string]Transaction)

	// Popola la mappa delle transazioni nella nuova catena
	for _, block := range newChain {
		for _, tx := range block.Transactions {
			// Usa una combinazione di sender, recipient, nonce e amount come chiave univoca
			key := fmt.Sprintf("%s:%s:%d:%f", tx.Sender, tx.Recipient, tx.Nonce, tx.Amount)
			txInNewChain[key] = true
		}
	}

	// Identifica le transazioni nella vecchia catena
	for _, block := range oldChain {
		for _, tx := range block.Transactions {
			// Ignora le transazioni di ricompensa
			if tx.Sender == "SYSTEM" {
				continue
			}
			key := fmt.Sprintf("%s:%s:%d:%f", tx.Sender, tx.Recipient, tx.Nonce, tx.Amount)
			txInOldChain[key] = tx
		}
	}

	// Identifica le transazioni orfane (presenti nella vecchia catena ma non nella nuova)
	var orphanTransactions []Transaction
	for key, tx := range txInOldChain {
		if !txInNewChain[key] {
			orphanTransactions = append(orphanTransactions, tx)
		}
	}

	return orphanTransactions
}

// Funzione per verificare la firma di una transazione
func VerifySignature(tx Transaction, pubkeyStr string) bool {
    // Decodifica la chiave pubblica
    pubkeyBytes, err := base64.StdEncoding.DecodeString(pubkeyStr)
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
    
    // Calcola l'hash della transazione (stesso metodo utilizzato per la firma)
    data := fmt.Sprintf("%s:%s:%f", tx.Sender, tx.Recipient, tx.Amount)
    hash := sha256.Sum256([]byte(data))
    
    // Decodifica la firma
    signatureBytes, err := base64.StdEncoding.DecodeString(tx.Signature)
    if err != nil {
        return false
    }
    
    // Estrai r e s dalla firma
    if len(signatureBytes) != 64 {
        return false
    }
    r := new(big.Int).SetBytes(signatureBytes[:32])
    s := new(big.Int).SetBytes(signatureBytes[32:])
    
    // Verifica la firma
    return ecdsa.Verify(pubkey, hash[:], r, s)
}

// Endpoint per ottenere il saldo di un indirizzo
http.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
    address := r.URL.Query().Get("address")
    if address == "" {
        http.Error(w, "Parametro 'address' mancante", http.StatusBadRequest)
        return
    }
    
    balance := blockchain.GetBalance(address)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "address": address,
        "balance": balance,
    })
})

// Endpoint per ottenere lo storico delle transazioni di un indirizzo
http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
    address := r.URL.Query().Get("address")
    if address == "" {
        http.Error(w, "Parametro 'address' mancante", http.StatusBadRequest)
        return
    }
    
    var transactions []Transaction
    
    // Cerca le transazioni in tutti i blocchi
    for _, block := range blockchain.Chain {
        for _, tx := range block.Transactions {
            if tx.Sender == address || tx.Recipient == address {
                transactions = append(transactions, tx)
            }
        }
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(transactions)
})

// Endpoint per ottenere tutti gli indirizzi e i loro saldi
http.HandleFunc("/addresses", func(w http.ResponseWriter, r *http.Request) {
    addresses := make(map[string]float64)
    
    // Trova tutti gli indirizzi unici
    for _, block := range blockchain.Chain {
        for _, tx := range block.Transactions {
            if tx.Sender != "SYSTEM" && tx.Sender != "" {
                addresses[tx.Sender] = 0
            }
            if tx.Recipient != "" {
                addresses[tx.Recipient] = 0
            }
        }
    }
    
    // Calcola il saldo per ogni indirizzo
    for address := range addresses {
        addresses[address] = blockchain.GetBalance(address)
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(addresses)
})

// Endpoint per validare la blockchain
http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
    isValid := blockchain.IsValid()
    
    w.Header().Set("Content-Type", "text/plain")
    if isValid {
        w.Write([]byte("La blockchain è valida!"))
    } else {
        w.Write([]byte("Errore: La blockchain non è valida!"))
    }
})

// Servi i file statici
http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, r.URL.Path[1:]) // Rimuovi lo slash iniziale
})

// Endpoint per l'explorer
http.HandleFunc("/explorer", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/explorer.html")
})

// Endpoint per la documentazione API
http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/api.html")
})

// Funzione principale
func main() {
	blockchain, err := LoadBlockchainFromFile("blockchain.json")
	if err != nil {
		blockchain = NewBlockchain(4) // Inizializza la blockchain se non esiste
	}

	// Inizializza il mempool vuoto
	pendingTransactions = []Transaction{}

	go syncWithPeers()

	// Pulizia periodica del mempool ogni minuto
	go func() {
		for {
			cleanMempool()
			time.Sleep(60 * time.Second)
		}
	}()

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
		// Controlla se la transazione è già presente nel mempool
		for _, pending := range pendingTransactions {
			if pending.Sender == tx.Sender && pending.Nonce == tx.Nonce && pending.Signature == tx.Signature {
				http.Error(w, "Transazione duplicata già presente nel mempool", http.StatusBadRequest)
				return
			}
		}
		// Controlla se la transazione è già presente nella blockchain
		for _, block := range blockchain.Chain {
			for _, txInBlock := range block.Transactions {
				if txInBlock.Sender == tx.Sender && txInBlock.Nonce == tx.Nonce && txInBlock.Signature == tx.Signature {
					http.Error(w, "Transazione duplicata già presente nella blockchain", http.StatusBadRequest)
					return
				}
			}
		}
		// Trova il nonce massimo già usato dal mittente (sia nella blockchain che nel mempool)
		maxNonce := -1
		for _, block := range blockchain.Chain {
			for _, txInBlock := range block.Transactions {
				if txInBlock.Sender == tx.Sender && txInBlock.Nonce > maxNonce {
					maxNonce = txInBlock.Nonce
				}
			}
		}
		for _, pending := range pendingTransactions {
			if pending.Sender == tx.Sender && pending.Nonce > maxNonce {
				maxNonce = pending.Nonce
			}
		}
		// Il nonce deve essere esattamente maxNonce+1
		if tx.Nonce != maxNonce+1 {
			http.Error(w, fmt.Sprintf("Nonce non valido: atteso %d, ricevuto %d", maxNonce+1, tx.Nonce), http.StatusBadRequest)
			return
		}
		// Controllo saldo solo se il mittente non è "SYSTEM"
		if tx.Sender != "SYSTEM" {
			// Calcola il saldo tenendo conto anche delle transazioni in sospeso
			saldo := blockchain.GetBalance(tx.Sender)
			for _, pending := range pendingTransactions {
				if pending.Sender == tx.Sender {
					saldo -= pending.Amount
				}
			}
			if saldo < tx.Amount {
				http.Error(w, "Saldo insufficiente per il mittente (considerando anche le transazioni in sospeso)", http.StatusBadRequest)
				return
			}
			if tx.Signature == "" {
				http.Error(w, "Transazione non firmata", http.StatusBadRequest)
				return
			}
		}
		// Sicurezza: importo positivo e mittente/destinatario diversi
		if tx.Amount < 0.0001 || tx.Amount > 1000000 {
			http.Error(w, "Importo non valido: troppo piccolo o troppo grande", http.StatusBadRequest)
			return
		}
		if math.IsNaN(tx.Amount) || math.IsInf(tx.Amount, 0) {
			http.Error(w, "Importo non valido", http.StatusBadRequest)
			return
		}
		if tx.Sender == tx.Recipient {
			http.Error(w, "Mittente e destinatario non possono essere uguali", http.StatusBadRequest)
			return
		}
		// Limita la lunghezza dei campi
		if len(tx.Sender) > 32 || len(tx.Recipient) > 32 || len(tx.Signature) > 128 {
			http.Error(w, "Uno dei campi supera la lunghezza massima consentita", http.StatusBadRequest)
			return
		}
		now := time.Now().Unix()
		if tx.Timestamp < now-3600 || tx.Timestamp > now+300 {
			http.Error(w, "Timestamp non valido: transazione troppo vecchia o nel futuro", http.StatusBadRequest)
			return
		}
		pendingTransactions = append(pendingTransactions, tx)
		cleanMempool()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Transazione aggiunta al mempool."))

		// Propaga la transazione ai peer
		go func(tx Transaction) {
			for _, peer := range peers {
				jsonTx, _ := json.Marshal(tx)
				http.Post(peer+"/receiveTransaction", "application/json", strings.NewReader(string(jsonTx)))
			}
		}(tx)
	})

	// Endpoint per minare tutte le transazioni in sospeso
	http.HandleFunc("/mine", func(w http.ResponseWriter, r *http.Request) {
		cleanMempool()
		if len(pendingTransactions) == 0 {
			http.Error(w, "Nessuna transazione da minare.", http.StatusBadRequest)
			return
		}
		// Ricevi il nome del miner dal body (JSON)
		var data struct {
			Miner string `json:"miner"`
		}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Dati di input non validi", http.StatusBadRequest)
			return
		}
		miner := data.Miner
		now := time.Now()
		if t, ok := lastMiningTime[miner]; ok && now.Sub(t) < miningCooldown {
			http.Error(w, fmt.Sprintf("Devi aspettare %v prima di minare di nuovo.", miningCooldown-now.Sub(t)), http.StatusTooManyRequests)
			return
		}
		lastMiningTime[miner] = now
		bc.MineBlock(pendingTransactions, miner)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Blocco minato con successo!"))
	})
}
