package chain

import (
	"bytes"
	"fmt"

	"github.com/boryoku-tekina/makiko/utils"
)

var (
	utxoPrefix   = []byte("utxo-")
	prefixLenght = len(utxoPrefix)
)

// GetUTXOOf function
// return all UTXO for the address
func GetUTXOOf(address string) TxOutputs {

	if IsInPool(address) {
		return UTXOInPoolOf(address)
	}

	UTXOs := TxOutputs{}
	UTXOs.Outputs = nil

	// wallets, err := wallet.CreateWallets()
	// utils.HandleErr(err)
	// w := wallets.GetWallets(address)

	PubKeyHash := utils.Base58Decode([]byte(address))
	PubKeyHash = PubKeyHash[1 : len(PubKeyHash)-4]

	lh := GetLastBlockHash()
	actualBlock := GetBlockByHash(lh)
Parcour:
	for {
		// if we are on the genesis block
		if bytes.Equal(actualBlock.Header.PrevHash, bytes.Repeat([]byte{0}, 32)) {
			break
		}
		if actualBlock.Transactions == nil {
			fmt.Println("there is no tx in this block")
		}
		// verifying Transactions in the block from the end to the start
		for i := len(actualBlock.Transactions) - 1; i >= 0; i-- {
			actualTx := actualBlock.Transactions[i]
			if actualTx.IsCoinBase() {
				for _, out := range actualTx.Outputs {
					if out.IsLockedWithKey(PubKeyHash) {
						fmt.Println("get coin base for ", address, " appending it...")
						UTXOs.Outputs = append(UTXOs.Outputs, out)
						fmt.Println("actual UTXOs.OUtputs : ", UTXOs.Outputs)
					}
				}
			} else {
				inOutput := false
				isChange := false
				var copy []TxOutput
				for _, out := range actualTx.Outputs {
					if out.IsLockedWithKey(PubKeyHash) {
						inOutput = true
						copy = append(copy, out)
					}
				}
				changeOutput := actualTx.Outputs[len(actualTx.Outputs)-1]

				isChange = changeOutput.IsLockedWithKey(PubKeyHash)

				if inOutput == true {
					if isChange == true {
						UTXOs.Outputs = append(UTXOs.Outputs, changeOutput)
						break Parcour
					} else {
						for _, out := range copy {
							UTXOs.Outputs = append(UTXOs.Outputs, out)
						}
					}
				}
			}

		}

		actualBlock = GetBlockByHash(actualBlock.Header.PrevHash)
	}
	return UTXOs
}

// if address contains in Inputs and OUtputs then reset the UTXO to the actual TxO
// if address contains only in Outpus then append actual UTXO with the actual TxO
