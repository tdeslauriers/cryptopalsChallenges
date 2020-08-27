package main

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

var key = []byte("YELLOW SUBMARINE")

func main() {

	file, err := ioutil.ReadFile("ch7.txt")
	if err != nil {
		panic(err.Error())
	}

	ct, err := base64.RawStdEncoding.DecodeString(string(file))
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("\nDecrypted text: \n\n%s", decryptAesEcb(ct, key))
}

func bytesDivisable(ct, k int) bool {

	if ct <= 0 {
		panic("Cipher text length is nil or zero bytes.")
	}

	if k <= 0 {
		panic("Key length is nil or zero bytes.")
	}

	if ct%k != 0 {
		return false
	}

	return true
}

func decryptAesEcb(ct, k []byte) []byte {

	block, err := aes.NewCipher(k)
	if err != nil {
		panic(err.Error())
	}

	if !bytesDivisable(len(ct), block.BlockSize()) {
		panic("Cipher text byte length not divisible by key byte length.")
	}

	for i := 0; i < len(ct); i += block.BlockSize() {

		block.Decrypt(ct[i:i+16], ct[i:i+16])
	}

	return ct
}
