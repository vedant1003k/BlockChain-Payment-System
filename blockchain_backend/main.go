package main

// [Code inserted here is from updated backend Go code in previous step]

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "strings"
    "sync"
    "time"
)

type Transaction struct {
    Sender    string  `json:"sender"`
    Receiver  string  `json:"receiver"`
    Amount    float64 `json:"amount"`
    Signature string  `json:"signature"`
}

type Block struct {
    Index        int           `json:"index"`
    Timestamp    string        `json:"timestamp"`
    Transactions []Transaction `json:"transactions"`
    PrevHash     string        `json:"prevHash"`
    Nonce        int           `json:"nonce"`
    Hash         string        `json:"hash"`
}

type Blockchain struct {
    Blocks          []Block       `json:"blocks"`
    TransactionPool []Transaction `json:"txPool"`
    mutex           sync.Mutex
}

var blockchain = Blockchain{
    Blocks: []Block{
        {
            Index:     0,
            Timestamp: time.Now().String(),
            Hash:      "0000",
            Nonce:     0,
        },
    },
}

func enableCORS(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins, or specify "http://localhost:3000"
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func calculateHash(block Block) string {
    record := strconv.Itoa(block.Index) + block.Timestamp + fmt.Sprintf("%v", block.Transactions) + block.PrevHash + strconv.Itoa(block.Nonce)
    h := sha256.New()
    h.Write([]byte(record))
    return hex.EncodeToString(h.Sum(nil))
}

func mineBlock(transactions []Transaction, prevBlock Block) Block {
    newBlock := Block{
        Index:        prevBlock.Index + 1,
        Timestamp:    time.Now().String(),
        Transactions: transactions,
        PrevHash:     prevBlock.Hash,
        Nonce:        0,
    }
    for {
        hash := calculateHash(newBlock)
        if strings.HasPrefix(hash, "0000") {
            newBlock.Hash = hash
            break
        }
        newBlock.Nonce++
    }
    return newBlock
}

func addBlock(newBlock Block) {
    blockchain.mutex.Lock()
    defer blockchain.mutex.Unlock()
    blockchain.Blocks = append(blockchain.Blocks, newBlock)
    blockchain.TransactionPool = []Transaction{}
}

func isValidTransaction(tx Transaction) bool {
    if tx.Sender == "MINT" {
        return true
    }
    balance := 0.0
    for _, block := range blockchain.Blocks {
        for _, t := range block.Transactions {
            if t.Sender == tx.Sender {
                balance -= t.Amount
            }
            if t.Receiver == tx.Sender {
                balance += t.Amount
            }
        }
    }
    return balance >= tx.Amount
}

func StartAPI() {
    http.HandleFunc("/wallet/new", func(w http.ResponseWriter, r *http.Request) {

        enableCORS(w) // Enable CORS for this route

        if r.Method == http.MethodOptions {
            return // Handle pre-flight OPTIONS request
        }

        priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
        pub := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
        json.NewEncoder(w).Encode(map[string]string{"public": hex.EncodeToString(pub)})
    })

    http.HandleFunc("/transaction", func(w http.ResponseWriter, r *http.Request) {

        enableCORS(w)
        var tx Transaction
        body, _ := ioutil.ReadAll(r.Body)
        json.Unmarshal(body, &tx)
        if !isValidTransaction(tx) {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("Insufficient balance or invalid transaction"))
            return
        }
        blockchain.TransactionPool = append(blockchain.TransactionPool, tx)
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"message": "Transaction queued"})
    })

    http.HandleFunc("/mine", func(w http.ResponseWriter, r *http.Request) {
        enableCORS(w)
        last := blockchain.Blocks[len(blockchain.Blocks)-1]
        newBlock := mineBlock(blockchain.TransactionPool, last)
        addBlock(newBlock)
        json.NewEncoder(w).Encode(newBlock)
    })

    http.HandleFunc("/blocks", func(w http.ResponseWriter, r *http.Request) {
        enableCORS(w)
        blockchain.mutex.Lock()
        defer blockchain.mutex.Unlock()
        json.NewEncoder(w).Encode(blockchain.Blocks)
    })

    http.HandleFunc("/balance/", func(w http.ResponseWriter, r *http.Request) {
        enableCORS(w)
        addr := strings.TrimPrefix(r.URL.Path, "/balance/")
        balance := 0.0
        var txs []Transaction
        for _, b := range blockchain.Blocks {
            for _, tx := range b.Transactions {
                if tx.Receiver == addr {
                    balance += tx.Amount
                    txs = append(txs, tx)
                }
                if tx.Sender == addr {
                    balance -= tx.Amount
                    txs = append(txs, tx)
                }
            }
        }
        json.NewEncoder(w).Encode(map[string]interface{}{
            "balance":      balance,
            "transactions": txs,
        })
    })

    fmt.Println("Server running at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
    StartAPI()
}
