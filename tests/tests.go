package tests

import (
	"fmt"
	"log"

	"github.com/Setra-Solofoniaina/glockchain/utils"

	"github.com/Setra-Solofoniaina/glockchain/chain"
	"github.com/Setra-Solofoniaina/glockchain/wallet"
)

// CreateWallet test the creation of wallet
func CreateWallet() {
	w, _ := wallet.CreateWallets()
	w.LoadFile()
	fmt.Println("creating wallet 3 times")
	w.AddWallet()
	w.AddWallet()
	w.AddWallet()
	fmt.Println("saving wallet file")
	w.SaveFile()
	all := w.GetAllAddresses()

	for _, addr := range all {
		fmt.Println(addr)
	}
}

// Transaction : test transactions functions
func Transaction() {
	chain.InitChain()
	cbtx := chain.CoinBaseTx(100, "1EFjGMygRKDw7dkFs2PLBFunhWpeFshKLp", "")
	chain.AddBlock([]*chain.Transaction{cbtx})
	tx := chain.NewTransaction("1EFjGMygRKDw7dkFs2PLBFunhWpeFshKLp", "1KgsVvVQzXgQ4LAaR3ZmFpDK9px4GD1ZXJ", 50)
	chain.AddBlock([]*chain.Transaction{tx})
}

// GetBalanceOf : get fund in address
func GetBalanceOf(address string) int {
	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	balance = chain.GetAmountOf(address)
	fmt.Printf("balance of %s == %d\n", address, balance)
	return balance
}

// Donate : give coin base transaction to an address
func Donate(address string, amount int) {
	T1 := chain.CoinBaseTx(amount, address, "Donation")
	var Txs []*chain.Transaction
	Txs = append(Txs, T1)
	chain.AddBlock(Txs)
	fmt.Printf("\n\n[INFO] : Donation of %d for %s DONE!\n\n", amount, address)
}

// Send : send amount of coin from an address to other
func Send(from, to string, amount int) {
	if from == to {
		log.Panic("DONT SEND COIN TO YOURSELF")
	}
	Tx := chain.NewTransaction(from, to, amount)
	var txs []*chain.Transaction
	txs = append(txs, Tx)
	chain.AddBlock(txs)
	fmt.Printf("[INFO]: Sending %d coins from %s to %s DONE!\n", amount, from, to)
}

// MemPoolTest Function to test mem Pool
func MemPoolTest() {
	chain.InitChain()
	Donate("1EFjGMygRKDw7dkFs2PLBFunhWpeFshKLp", 150)
	Donate("1KgsVvVQzXgQ4LAaR3ZmFpDK9px4GD1ZXJ", 1)
	chain.Send("1EFjGMygRKDw7dkFs2PLBFunhWpeFshKLp", "1KgsVvVQzXgQ4LAaR3ZmFpDK9px4GD1ZXJ", 50)
	chain.MinePendingTx()
	chain.Send("1EFjGMygRKDw7dkFs2PLBFunhWpeFshKLp", "1KgsVvVQzXgQ4LAaR3ZmFpDK9px4GD1ZXJ", 100)
	fmt.Println("last valid tx added to mempool")
	chain.Send("1KgsVvVQzXgQ4LAaR3ZmFpDK9px4GD1ZXJ", "1EFjGMygRKDw7dkFs2PLBFunhWpeFshKLp", 100)
	fmt.Println("we must hit an exception")
}

// MerkleTreeTest Function to test merkleTree function
func MerkleTreeTest() {
	chain.InitChain()
	Donate("1EFjGMygRKDw7dkFs2PLBFunhWpeFshKLp", 150)
	Donate("1KgsVvVQzXgQ4LAaR3ZmFpDK9px4GD1ZXJ", 1)
	chain.Send("1EFjGMygRKDw7dkFs2PLBFunhWpeFshKLp", "1KgsVvVQzXgQ4LAaR3ZmFpDK9px4GD1ZXJ", 50)
	chain.MinePendingTx()
	b1 := chain.GetBlockByHash(chain.GetLastBlockHash())
	fmt.Printf("HMR : %x\n", b1.Header.HMerkleRoot)
	chain.Send("1EFjGMygRKDw7dkFs2PLBFunhWpeFshKLp", "1KgsVvVQzXgQ4LAaR3ZmFpDK9px4GD1ZXJ", 100)
	chain.MinePendingTx()
	b := chain.GetBlockByHash(chain.GetLastBlockHash())
	fmt.Printf("HMR : %x\n", b.Header.HMerkleRoot)
	fmt.Println("voila")

}
