package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/boryoku-tekina/makiko/utils"
	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

// Wallet represent an User Account
// ecdsa : elliptic curve digital signing algorithm
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewKeyPair create a key pair (private and public) for a wallet
// return a private key and a public corresponding it
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	// curve we use to generate the private key
	curve := elliptic.P256() // the output of a curve will be 256 byte

	private, err := ecdsa.GenerateKey(curve, rand.Reader)

	utils.HandleErr(err)

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub

}

// MakeWallet create wallet from two keys(private and public)
func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

/*
	[--------------STEP TO HAVE THE ADDRESS-----------------]

						[private key]
								!
							[ecdsa]
								!
						[public key hash]
								!
							[sha256 key]
								!
							[ripemd160]
								!
						[public key hash]---------------+
								!                      	!
							[sha256]                   	!
								!                      	!
							[sha256]                   	!
								!                      	!
							[4bytes]                   	!
								!                      	!
							[checksum]                 	!                    [version]
								!                      	!                        !
								+------------------>[base58]<--------------------+
														!
														!
													[ADDRESS]
*/

// PublicKeyHash return the hash of the public key
func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])

	utils.HandleErr(err)

	publicRipMD := hasher.Sum(nil)

	return publicRipMD
}

// CheckSum : get the first 4bytes checksum of the hash
func CheckSum(payload []byte) []byte {
	// hashing it 2 times
	fistHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(fistHash[:])

	// return the first 4byte of the hash
	return secondHash[:checksumLength]
}

// Address function to get the final address
// return the final address to receive coins
func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)
	checksum := CheckSum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	finalAddress := utils.Base58Encode(fullHash)

	// fmt.Printf("pub key :\t%x \n", w.PublicKey)
	// fmt.Printf("pub Hash :\t%x \n", pubHash)
	// fmt.Printf("vHash : \t%x\n", versionedHash)
	// fmt.Printf("checksum : \t%x\n", checksum)
	// fmt.Printf("fullHash : \t%x\n", fullHash)
	// fmt.Printf("address :\t%s \n", finalAddress)

	return finalAddress
}

// ValidateAddress return true if address is valid
func ValidateAddress(address string) bool {
	pubKeyHash := utils.Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	targetChecksum := CheckSum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}
