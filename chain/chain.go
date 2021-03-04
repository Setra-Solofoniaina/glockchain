package chain

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boryoku-tekina/makiko/utils"
	"github.com/boryoku-tekina/makiko/wallet"
)

// InitChain initialize the blockchain
// create the database and the first transaction ── Genesis
// store the last hash key/value to the Genesis db
// One block ──> One Database = one file ; so we can download one DB file instead of
// downloading all the database from the beginning and file by file
func InitChain() {
	if !DBExists() {
		fmt.Println("[INFO] : no chain yet, creating genesis block")
		CreateGenesisBlock()
	}
	fmt.Println("[INFO] : it means that there is already a chain database")
}

// AddBlock : Add New Block to the chain
// mine the block with the pending txs
func AddBlock(transactions []*Transaction) {
	var b Block
	b.Transactions = transactions
	b.Mine()
}

// PrintChain : print the chain
func PrintChain() {
	// iterating through all block
	// beginning from the last
	lh := GetLastBlockHash()
	actualBlock := GetBlockByHash(lh)
	for {
		// if we are on the genesis block
		if bytes.Equal(actualBlock.Header.PrevHash, bytes.Repeat([]byte{0}, 32)) {
			actualBlock.PrintBlockInfo()
			break
		}
		actualBlock.PrintBlockInfo()
		actualBlock = GetBlockByHash(actualBlock.Header.PrevHash)
	}
}

// ValidChain return true if chain is valid
// if all block is connected
func ValidChain() bool {
	// var actualBlock Block
	lh := GetLastBlockHash()
	actualBlock := GetBlockByHash(lh)

	for {
		// if we are on the genesis block
		if bytes.Equal(actualBlock.Header.PrevHash, bytes.Repeat([]byte{0}, 32)) {
			validation := actualBlock.ValidateBlock()
			if validation == false {
				return false
			}
			break
		}
		validation := actualBlock.ValidateBlock()
		if validation == false {
			return false
		}
		actualBlock = GetBlockByHash(actualBlock.Header.PrevHash)
	}
	return true
}

// DBExists function to check if database file exist
func DBExists() bool {
	if _, err := os.Stat("./DB/LastBlockHash.bc"); os.IsNotExist(err) {
		return false
	}

	return true
}

// SignTransaction function to sign a tx
func SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	// prevTxs := make(map[string]Transaction)

	// for _, in := range tx.Inputs {
	// 	prevTX, err := FindTransaction(in.ID)
	// 	utils.HandleErr(err)
	// 	prevTxs[hex.EncodeToString(prevTX.ID)] = prevTX
	// }
	tx.Sign(privKey)
}

// FindTransaction : find transaction by given id
func FindTransaction(ID []byte) (Transaction, error) {
	lh := GetLastBlockHash()
	actualBlock := GetBlockByHash(lh)

	for {
		// break if we are on the genesis block
		if bytes.Equal(actualBlock.Header.PrevHash, bytes.Repeat([]byte{0}, 32)) {
			break
		}

		for _, tx := range actualBlock.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
		actualBlock = GetBlockByHash(actualBlock.Header.PrevHash)
	}

	return Transaction{}, errors.New("Transaction does not exist")
}

// NewTransaction : create new transaction from an address to another adress
func NewTransaction(from, to string, amount int) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	wallets, err := wallet.CreateWallets()
	utils.HandleErr(err)
	w := wallets.GetWallets(from)
	// pubKeyHash := wallet.PublicKeyHash(w.PublicKey)

	fromBalance := GetAmountOf(from)
	fmt.Printf("amount of %s : %d\n", from, fromBalance)
	if fromBalance < amount {
		log.Panic("[ERROR] : Not enough coins")
	}

	UTXO := GetUTXOOf(from)

	fmt.Println("[UTXO] : UTXOS")
	fmt.Println(UTXO)

	for _, outs := range UTXO.Outputs {
		input := TxInput{outs.Value, nil, w.PublicKey}
		inputs = append(inputs, input)
	}

	outputs = append(outputs, *NewTxOutput(amount, to))
	// CHANGES
	fmt.Printf("[INFO] : FROM BALANCE : %d\n", fromBalance)
	fmt.Printf("[INFO] : AMOUNT TO SEND : %d\n", amount)
	change := fromBalance - amount
	fmt.Printf("[INFO] : CHANGES : %d\n", change)
	outputs = append(outputs, *NewTxOutput(change, from))

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	SignTransaction(&tx, w.PrivateKey)

	return &tx

}

// Send : send amount of coin from an address to other
func Send(from, to string, amount int) {
	if from == to {
		log.Panic("DONT SEND COIN TO YOURSELF")
	}
	Tx := NewTransaction(from, to, amount)
	Tx.AddToPool()
	fmt.Printf("[INFO] : Transaction added to mem pool")
}
