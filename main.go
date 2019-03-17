package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Block struct {
	Index     int
	Timestamp string
	BPM       int
	PrevHash  string
	Hash      string
}

var Blockchain []Block

type Message struct {
	BPM int
}

var mutex = &sync.Mutex{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := Block{0, t.String(), 0, "", ""}
		spew.Dump(genesisBlock)

		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		mutex.Unlock()
	}()
	log.Fatal(run())

}

func run() error {
	r := setupRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", httpAddr)
	r.Run("localhost:" + httpAddr)
	return nil
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", handleGetBlockchain)
	r.POST("/", handleWriteBlock)
	return r
}

func handleGetBlockchain(c *gin.Context) {
	c.JSON(http.StatusOK, Blockchain)
}

func handleWriteBlock(c *gin.Context) {
	var msg Message
	c.BindJSON(&msg)

	mutex.Lock()
	prevBlock := Blockchain[len(Blockchain)-1]
	newBlock := generateBlock(prevBlock, msg.BPM)

	if isBlockValid(newBlock, prevBlock) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
		spew.Dump(Blockchain)
	}
	mutex.Unlock()

	c.JSON(http.StatusOK, newBlock)
}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, BPM int) Block {
	var newBlock Block
	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}
