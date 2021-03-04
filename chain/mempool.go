package chain

import (
	"github.com/boryoku-tekina/makiko/utils"
)

// pool contain all unmined transactions
var pool []Transaction

// GetPool return Transaction in mempool
func GetPool() []Transaction {
	return pool
}

// AddToPool add a transaction into pool memory
func (t *Transaction) AddToPool() {
	pool = append(pool, *t)
}

// CountTxInPool return number of transaction in memory pool
func CountTxInPool() int {
	return len(pool) - 1
}

// IsInPool : return true if an address is in Inputs
// in Transactions in mempool
func IsInPool(address string) bool {
	PubKeyHash := utils.Base58Decode([]byte(address))
	PubKeyHash = PubKeyHash[1 : len(PubKeyHash)-4]

	for _, tx := range pool {
		for _, in := range tx.Inputs {
			if in.UsesKey(PubKeyHash) {
				return true
			}
		}
	}
	return false
}

// UTXOInPoolOf : return the UTXO of an address in memPool
func UTXOInPoolOf(address string) TxOutputs {
	UTXOs := TxOutputs{}
	UTXOs.Outputs = nil
	isChange := false
	PubKeyHash := GetPubKeyHash(address)
Parcour:
	for i := len(pool) - 1; i >= 0; i-- {
		changeOutput := pool[i].Outputs[len(pool[i].Outputs)-1]

		isChange = changeOutput.IsLockedWithKey(PubKeyHash)

		if isChange == true {
			UTXOs.Outputs = append(UTXOs.Outputs, changeOutput)
			break Parcour
		}
	}
	return UTXOs
}

// purgePool : reset the mempool
func purgePool() {
	pool = []Transaction{}
}

// MinePendingTx : mine all tx in memPool
func MinePendingTx() {
	tomine := pool
	purgePool()
	var txs []*Transaction
	for _, p := range tomine {
		txs = append(txs, &p)
	}
	AddBlock(txs)
}
