package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-ipfs-api"
	"github.com/joho/godotenv"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

var h hash.Hash
var aesgcm cipher.AEAD
var nonce []byte

var mutex = &sync.Mutex{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

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
	secret := os.Getenv("SECRET")
	key := []byte(secret)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	nonce = make([]byte, 12) // TODO: How large should the nonce be?
	if false {
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Nonce: %x\n", nonce)

	aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		log.Fatalln(err)
	}

	h = sha256.New()
	sh = shell.NewShell("localhost:5001")

	// Validate that connection is active
	if _, err := sh.ID(); err != nil {
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
	c.JSON(http.StatusOK, Blockchain)
}

func handleGetBlockData(c *gin.Context) {
	cid := c.Params.ByName("cid")

	objBytes, err := fetchObjectFromIPFS(cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	plaintext, err := aesgcm.Open(nil, nonce, objBytes, nil)
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

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.IPFSHash) + block.PrevHash
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block
	t := time.Now()

	var err error
	BPMuint := uint32(BPM)
	plaintext := make([]byte, 4)
	binary.BigEndian.PutUint32(plaintext, BPMuint)

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
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
