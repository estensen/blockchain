package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/estensen/blockchain/crypto"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-ipfs-api"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Block struct {
	Index     int
	Timestamp string
	IPFSHash  string
	PrevHash  string
	Hash      string
}

var Blockchain []Block

type Message struct {
	BPM int
}

var sh *shell.Shell

var cryptor *crypto.Cryptor

var mutex = &sync.Mutex{}

func main() {
	go func() {
		t := time.Now()
		genesisBlock := Block{0, t.String(), "", "", ""}
		spew.Dump(genesisBlock)

		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		mutex.Unlock()
	}()

	log.Fatal(run())
}

func run() error {
	var err error

	if err := godotenv.Load(); err != nil {
		return err
	}

	userPass := os.Getenv("SECRET")
	if userPass == "" {
		userPass = crypto.GetEncPass()
	}

	if cryptor, err = crypto.NewCryptor(userPass); err != nil {
		return err
	}

	sh = shell.NewShell("localhost:5001")
	// Validate that connection is active
	if _, err := sh.ID(); err != nil {
		fmt.Println("You are probably not running the IPFS daemon")
		return err
	}

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
	r.GET("/block/:cid", handleGetBlockData)
	return r
}

func handleGetBlockchain(c *gin.Context) {
	fmt.Println("The entire blockchain")
	spew.Dump(Blockchain)
	c.JSON(http.StatusOK, Blockchain)
}

func handleGetBlockData(c *gin.Context) {
	cid := c.Params.ByName("cid")

	data, err := fetchObjectFromIPFS(cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	plaintext, err := cryptor.Decrypt(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	BPM := binary.BigEndian.Uint32(plaintext)

	c.JSON(http.StatusOK, gin.H{"BPM": BPM})
}

func handleWriteBlock(c *gin.Context) {
	var msg Message
	c.BindJSON(&msg)

	mutex.Lock()
	prevBlock := Blockchain[len(Blockchain)-1]
	newBlock, err := generateBlock(prevBlock, msg.BPM)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else if isBlockValid(newBlock, prevBlock) {
		newBlockchain := append(Blockchain, newBlock)
		fmt.Println("The new b")
		spew.Dump(newBlockchain)
		replaceChain(newBlockchain)
		spew.Dump(Blockchain)
	}
	mutex.Unlock()

	c.JSON(http.StatusOK, newBlock)
}

func fetchObjectFromIPFS(cid string) ([]byte, error) {
	r, err := sh.Cat(cid)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	objBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return objBytes, nil
}

func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + block.IPFSHash + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block
	t := time.Now()

	var err error
	// TODO: Use variable-length encoding
	BPMuint := uint32(BPM)
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, BPMuint)

	ciphertext := cryptor.Encrypt(data)
	fmt.Printf("Ciphertext: %x\n", ciphertext)

	newBlock.IPFSHash, err = sh.Add(bytes.NewReader(ciphertext))
	if err != nil {
		log.Fatal(err)
	}

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, err
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		fmt.Println("The block number is not incremented by one")
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		fmt.Println("The hash of the previous block is not referenced correctly")
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		fmt.Println("The block hash is not correct")
		return false
	}

	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}
