package main

import (
	"encoding/hex"
	"fmt"
	"log"
)

// omitted testing because too simple
func main() {

	s1 := "1c0111001f010100061a024b53535009181c"
	s2 := "686974207468652062756c6c277320657965"

	// should be 746865206b696420646f6e277420706c6179
	r := xor(hexToBytes(s1), hexToBytes(s2))
	fmt.Printf("Resulting hex: %x\nand the same in plain text: %q\n", r, r)
}

func hexToBytes(s string) (b []byte) {

	str := []byte(s)
	b = make([]byte, hex.DecodedLen(len(str)))
	hex.Decode(b, str)

	return
}

func xor(b1, b2 []byte) (r []byte) {
	if len(b1) != len(b2) {
		log.Panic("Input lengths do not match.")
		return nil
	}

	for i := range b1 {
		r = append(r, b1[i]^b2[i])
	}

	return
}
