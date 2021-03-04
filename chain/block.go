package chain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/boryoku-tekina/makiko/utils"
)

const (
	dBPath = "DB/"
)

// BHeader represent block Header
type BHeader struct {
	Hash        []byte
	HMerkleRoot []byte
	PrevHash    []byte
	Timestamp   []byte
	Nonce       int
}

// Block represent one block structure
type Block struct {
	Header       BHeader
	Transactions []*Transaction
}

// OpenDatabase : opening a database
func OpenDatabase(s string) *bolt.DB {
	opts := bolt.DefaultOptions

	Path := dBPath + s + ".bc"

	db, err := bolt.Open(Path, 0777, opts)
	utils.HandleErr(err)
	return db
}

// CreateGenesisBlock : create the genesis Block and mine it
func CreateGenesisBlock() {
	fmt.Println("[INFO] : Creating Genesis block")
	var b Block
	b.Transactions = nil
	b.Header.Timestamp = []byte(time.Now().String())
	b.Header.PrevHash = bytes.Repeat([]byte{0}, 32)

	pow := NewWork(&b)
	nonce, hash := pow.Work()

	b.Header.Hash = hash[:]
	b.Header.Nonce = nonce

	if !b.ValidateBlock() {
		log.Panic("this block is not valid")
	}
	b.RegisterToDB()
	b.SetAsLastBlock()

	fmt.Println("Genesis Block Generated")
}

// Mine : create a new block and mine it
func (b *Block) Mine() {
	// block := &Block{[]byte{}, txs, prevHash, 0}

	b.Header.Timestamp = []byte(time.Now().String())
	b.SetHMerkleRoot()
	b.Header.PrevHash = GetLastBlockHash()
	// GetLastBlockHash()

	pow := NewWork(b)
	nonce, hash := pow.Work()

	b.Header.Hash = hash[:]
	b.Header.Nonce = nonce

	if !b.ValidateBlock() {
		log.Panic("[ERROR] : this block is not valid")
	}
	b.RegisterToDB()
	b.SetAsLastBlock()
	fmt.Println("[INFO] : Block Mined Successfully")
}

// ValidateBlock : validate a passed block
func (b *Block) ValidateBlock() bool {
	pow := NewWork(b)
	return pow.Validate()
}

// RegisterToDB open new boltDB database and create the block value one key in it
func (b *Block) RegisterToDB() {

	db := OpenDatabase(fmt.Sprintf("%x", b.Header.Hash))
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(b.Header.Hash)
		utils.HandleErr(err)
		bucket.Put([]byte("block"), b.Serialize())
		return err
	})
	utils.HandleErr(err)
	fmt.Println("block successfully registered to database")
}

// SetAsLastBlock set a block as the last block in the database
func (b *Block) SetAsLastBlock() {

	db := OpenDatabase("LastBlockHash")
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("LastBlockHash"))
		utils.HandleErr(err)
		bucket.Put([]byte("LastBlockHash"), b.Header.Hash)
		return err
	})
	utils.HandleErr(err)

	fmt.Println("[INFO] : block successfully set as last block")
}

// GetLastBlockHash get the hash of the last block in the chain
func GetLastBlockHash() []byte {
	var lh []byte
	getLastBlockHash(&lh)
	return lh
}

// GetLastBlockHash get the hash of the last block in the chain
// putting it in d
func getLastBlockHash(d *[]byte) {
	// var result []byte

	var b Block

	db := OpenDatabase("LastBlockHash")
	defer db.Close()

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("LastBlockHash"))
		resultat := bucket.Get([]byte("LastBlockHash"))
		b.Header.Hash = resultat

		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		errs := encoder.Encode(resultat)
		utils.HandleErr(errs)
		res := buffer.Bytes()
		// fmt.Printf("Inspecting res : %x\n", res[0:4])
		*d = res[4:]

		return nil
	})
	utils.HandleErr(err)
	fmt.Println("[INFO] : got last block hash")
}

// GetBlockByHash : return a block according to the given hash
func GetBlockByHash(bhash []byte) *Block {
	db := OpenDatabase(fmt.Sprintf("%x", bhash))
	defer db.Close()
	var b *Block

	// var pb *Block
	var byterep *[]byte

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bhash)
		b = Deserialize(bucket.Get([]byte("block")))

		bBytes := b.Serialize()

		byterep = &bBytes

		return nil
	})

	utils.HandleErr(err)

	return Deserialize(*byterep)
}

// PrintBlockInfo : print all information in the block
func (b *Block) PrintBlockInfo() {
	fmt.Println("────────────────────────────────────────")
	fmt.Printf("Block  %x information : \n", b.Header.Hash)

	fmt.Printf("────Transactions :\n\t")
	for _, tx := range b.Transactions {
		fmt.Printf("----------------------------------------------------------\n")
		fmt.Printf(tx.String())
		fmt.Printf("\n----------------------------------------------------------")
	}
	fmt.Println()
	fmt.Printf("────Hash : \t %x \n", b.Header.Hash)
	fmt.Printf("────Previous Hash : \t %x \n", b.Header.PrevHash)
	fmt.Printf("────Timestamp : \t %s \n", string(b.Header.Timestamp))
	fmt.Printf("────Nonce : \t %d \n", b.Header.Nonce)
	fmt.Println("────────────────────────────────────────")

}

// UTILS FUNCTIONS

// Serialize [return a byte representation of a BLOCK]
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	utils.HandleErr(err)

	return res.Bytes()
}

// Deserialize a block
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("error while decoding")
	}
	utils.HandleErr(err)

	return &block
}
