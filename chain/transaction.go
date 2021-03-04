package chain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/boryoku-tekina/makiko/utils"
)

var minerReward int = 100

// Transaction : represent a transaction
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// String function
func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("──Transaction id :  %x:", tx.ID))
	lines = append(lines, fmt.Sprintf("----------------------------------INPUTS--------------------------------------------"))

	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("\t Input:\t %d", i))
		lines = append(lines, fmt.Sprintf("\t Out:\t %d", input.Out))
		lines = append(lines, fmt.Sprintf("\t Signature:\t %x", input.Signature))
		lines = append(lines, fmt.Sprintf("\t PubKey:\t %x", input.PubKey))
		lines = append(lines, fmt.Sprintf("----------------------------------"))

	}
	lines = append(lines, fmt.Sprintf("----------------------------------OUTPUTS--------------------------------------------"))

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("\t Output:\t %d", i))
		lines = append(lines, fmt.Sprintf("\t Value:\t %d", output.Value))
		lines = append(lines, fmt.Sprintf("\t Script:\t %x", output.PubKeyHash))

		lines = append(lines, fmt.Sprintf("----------------------------------"))

	}

	return strings.Join(lines, "\n")

}

// SetID : setting id
// the tx id is the hash of the encoded tx
// (encoded tx = byte representation of tx)
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	utils.HandleErr(err)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// Serialize : return a []byte representation of a transaction
func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	utils.HandleErr(err)

	return encoded.Bytes()
}

// Hash : Hash the transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// IsCoinBase FUnction
// Determing what type of tx we have
func (tx *Transaction) IsCoinBase() bool {
	// coin base only have ONE input
	// return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
	return len(tx.Inputs) == 0
}

// CoinBaseTx : transaction of a coin
func CoinBaseTx(amount int, to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("New generated Coin to %s", to)
	}
	// txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTxOutput(amount, to)

	// tx := Transaction{nil, []TxInput{txin}, []TxOutput{*txout}}
	tx := Transaction{nil, nil, []TxOutput{*txout}}

	tx.SetID()

	return &tx
}

// Sign : Sign a transaction
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey) {
	if tx.IsCoinBase() {
		return
	}
	tx.ID = tx.Hash()
	r, s, err := ecdsa.Sign(rand.Reader, &privKey, tx.ID)
	utils.HandleErr(err)
	signature := append(r.Bytes(), s.Bytes()...)

	for _, in := range tx.Inputs {
		in.Signature = signature
	}
}

// GetAmountOf Address
func GetAmountOf(address string) int {
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	amount := 0

	UTXOs := GetUTXOOf(address)

	for _, utx := range UTXOs.Outputs {
		if utx.IsLockedWithKey(pubKeyHash) {
			amount += utx.Value
		}
	}
	return amount
}

// CountTransaction return number of transaction in a block
func (b *Block) CountTransaction() int {
	return len(b.Transactions)
}
