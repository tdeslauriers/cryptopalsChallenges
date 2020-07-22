package main

import (
	"fmt"
	"os"
)

func main() {

	t := []byte(os.Args[1])
	k := []byte(os.Args[2])
	x := repeatXor(t, k)

	fmt.Printf("Xor encoded hex: %x\n", x)
	fmt.Printf("Xor decoded text: %s\n", repeatXor(x, k))
}

func repeatXor(t, k []byte) (c []byte) {

	i := 0
	for j := 0; j < len(t); j++ {

		c = append(c, t[j]^k[i])
		if i == len(k)-1 {
			i = 0
		} else {
			i++
		}
	}

	return
}
