package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"syscall"
)

// TODO: Increase after testing
const passLen = 4

// Generated with: head -c 32 /dev/urandom | base64
// TODO: Get new salt at runtime
const salt = "5MieAyX5FLxVHU4CFpMiHyz8v3O//vAyHbP2xQVTwos="

type Cryptor struct {
	Aead cipher.AEAD
}

func (c *Cryptor) Encrypt(data []byte) []byte {
	nonceSize := c.Aead.NonceSize()
	nonce := make([]byte, nonceSize)
	io.ReadFull(rand.Reader, nonce)
	return c.Aead.Seal(nonce, nonce, data, nil)
}

func (c *Cryptor) Decrypt(data []byte) ([]byte, error) {
	nonceSize := c.Aead.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return c.Aead.Open(nil, nonce, ciphertext, nil)
}

func NewCryptor(passphrase string) (*Cryptor, error) {
	c := new(Cryptor)

	kdf := argon2.Key([]byte(passphrase), []byte(salt), 3, 32*1024, 4, 44)

	block, err := aes.NewCipher(kdf[:32])
	if err != nil {
		return c, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return c, err
	}

	c.Aead = aead

	return c, nil
}

func GetEncPass() string {
	pass1 := getPass("Enter password: ")
	pass2 := getPass("Retype password: ")

	if len(pass1) < passLen {
		fmt.Printf("The password must be at least %d characters\n", passLen)
		os.Exit(0)
	}

	if pass1 != pass2 {
		fmt.Println("The passwords do not match")
		os.Exit(0)
	}

	return pass1
}

// Get user password without echoing to the screen
func getPass(promt string) string {
	fmt.Print(promt)
	pass, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Could not read password")
		os.Exit(0)
	}

	fmt.Println("")

	return string(pass)
}
