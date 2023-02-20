package proofcrypto

import (
	"crypto/sha256"
	"fmt"
	"github.com/multiformats/go-multihash"
	"hash/fnv"
	"math/rand"
	"time"
)

// HashFile
// Hash the file contents
func HashFile(fileContents string) string {
	hash := sha256.Sum256([]byte(fileContents))
	fmt.Printf("Hash of modified file contents: %x\n", hash)

	return fmt.Sprintf("%x", hash)
}

// CreateRandomHash
// Create a random hash
func CreateRandomHash() string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	randomNumber := random.Intn(5000000000) + 1
	hash, _ := multihash.Sum([]byte(fmt.Sprintf("%d", randomNumber)), multihash.SHA2_256, -1)
	return hash.B58String()
}

// GetIntFromHash
// get a random number from a hash
func GetIntFromHash(hash string) uint32 {
	h := fnv.New32a()
	// Write the input string to the hash object
	h.Write([]byte(hash))
	// Get the 32-bit hash value as a uint32 and convert to a number between 1 and 100
	hashValue := h.Sum32()%10 + 1
	return hashValue
}
