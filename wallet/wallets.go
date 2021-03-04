package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/boryoku-tekina/makiko/utils"
)

const walletFile = "./tmp/wallets.data"

// Wallets represent saved wallet of the use
type Wallets struct {
	Wallets map[string]*Wallet
}

// CreateWallets create and load the wallets file
func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile()
	return &wallets, err
}

// AddWallet : add new address wallet
func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := fmt.Sprintf("%s", wallet.Address())

	ws.Wallets[address] = wallet

	return address
}

// GetAllAddresses allows to get all of the addresses according to the wallet
func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

// GetWallets function
func (ws Wallets) GetWallets(address string) Wallet {
	return *ws.Wallets[address]
}

// LoadFile load the saved wallets file
func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	var wallets Wallets

	fileContent, err := ioutil.ReadFile(walletFile)

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)

	utils.HandleErr(err)

	ws.Wallets = wallets.Wallets
	// we have no error so return nil
	return nil
}

// SaveFile save the wallet file
func (ws *Wallets) SaveFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)

	utils.HandleErr(err)

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	utils.HandleErr(err)
}
