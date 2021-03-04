package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"

	"github.com/mr-tron/base58"
)

// HandleErr : Handle error function
func HandleErr(err error) {
	if err != nil {
		log.Panic("[ERROR] : ", err.Error())
	}
}

// ToHex function
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// Base58Encode encode to a base58
func Base58Encode(input []byte) []byte {
	encoded := base58.Encode(input)

	return []byte(encoded)
}

// Base58Decode decode a base58
func Base58Decode(input []byte) []byte {
	decoded, err := base58.Decode(string(input[:]))

	HandleErr(err)

	return decoded

}

// CreateDBFolder create Database folder if not exist
func CreateDBFolder() {
	_, err := os.Stat("DB")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("DB", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
}
