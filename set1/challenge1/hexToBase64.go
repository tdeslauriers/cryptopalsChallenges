package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

var input = []byte("49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d")

func main() {

	p, err := DecodeHex(input)
	if err != nil {
		fmt.Printf("Failed to decode hex: %s", err)
	}
	e := Base64Encode(p)
	fmt.Printf("Base64 outcome: %s\n", e)

}

func DecodeHex(input []byte) ([]byte, error) {
	pt := make([]byte, hex.DecodedLen(len(input)))
	_, err := hex.Decode(pt, input)
	if err != nil {
		return nil, err
	}
	return pt, nil
}

func Base64Encode(input []byte) []byte {
	e := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
	base64.StdEncoding.Encode(e, input)

	return e
}
