package chain

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/boryoku-tekina/makiko/utils"
)

// take the data from the block

// create a counter (nonce) which start at 0

// create a hash of the data plus the counter

// check the hash to see if it meet a set of requirements

// REQUIREMENTS:
// THe first few bytes of the hash must contain 0

// Difficulty of consensus
var difficulty = 1

// ProofOfWork struct
type ProofOfWork struct {
	Block *Block
}

// NewWork return a new pow
func NewWork(b *Block) *ProofOfWork {
	pow := &ProofOfWork{b}
	return pow
}

// initData Function
func (pow *ProofOfWork) initData(nonce int) []byte {

	data := bytes.Join(
		[][]byte{
			pow.Block.Header.PrevHash,
			pow.Block.Header.HMerkleRoot,
			[]byte(pow.Block.Header.Timestamp),
			utils.ToHex(int64(nonce)),
			utils.ToHex(int64(difficulty)),
		},
		[]byte{},
	)
	return data
}

// Work the PoW algorithm
func (pow *ProofOfWork) Work() (int, []byte) {
	var hash [32]byte
	nonce := 0

	for {
		data := pow.initData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)

		if bytes.HasPrefix(hash[:], bytes.Repeat([]byte{0}, difficulty)) {
			break
		}
		nonce++
	}
	fmt.Println()
	return nonce, hash[:]
}

// Validate the work
func (pow *ProofOfWork) Validate() bool {
	data := pow.initData(pow.Block.Header.Nonce)

	hash := sha256.Sum256(data)

	if !bytes.HasPrefix(hash[:], bytes.Repeat([]byte{0}, difficulty)) {
		fmt.Println("hash does not satisfy difficulty requirements")
		return false
	}
	if !bytes.Equal(hash[:], pow.Block.Header.Hash) {
		fmt.Println("hash is not correct : maybe this nonce does not provide a valid hash")
		return false
	}
	// verifying the genesis block
	if bytes.Equal(pow.Block.Header.PrevHash, bytes.Repeat([]byte{0}, 32)) {
		fmt.Println("GENESIS BLOCK VALID")
		return true
	}
	if !bytes.Equal(pow.Block.Header.PrevHash, GetBlockByHash(pow.Block.Header.PrevHash).Header.Hash) {
		fmt.Println("previous hash not matching")
	}
	fmt.Printf("[INFO] : block %x validated\n", pow.Block.Header.Hash)
	return true
}
