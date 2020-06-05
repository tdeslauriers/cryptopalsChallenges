package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
)

var s1 = "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
var keys = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// http://pi.math.cornell.edu/~mec/2003-2004/cryptography/subs/frequencies.html
var Fq = map[byte]float64{'E': 12.02, 'T': 9.10, 'A': 8.12, 'O': 7.68, 'I': 7.31, 'N': 6.95, 'S': 6.28, 'R': 6.02, 'H': 5.92, 'D': 4.32, 'L': 3.98, 'U': 2.88, 'C': 2.71, 'M': 2.61, 'F': 2.30, 'Y': 2.11, 'W': 2.09, 'G': 2.03, 'P': 1.82, 'B': 1.49, 'V': 1.11, 'K': 0.69, 'X': 0.17, 'Q': 0.11, 'J': 0.10, 'Z': 0.07}

type cipher struct {
	key     byte
	message []byte
	freq    float64
}

type ciphers []cipher

var wg sync.WaitGroup

func main() {

	pt := hexToBytes(s1)
	k := []byte(keys)

	wg.Add(len(keys))
	c := make(chan cipher)

	for _, j := range k {
		go func(b byte) {
			dt := xor(pt, b)
			fq := calcFq(dt)
			ci := cipher{
				key:     b,
				message: dt,
				freq:    fq,
			}
			c <- ci
			wg.Done()
		}(j)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	var ciphs ciphers

	for n := range c {
		ciphs = append(ciphs, n)
	}

	sort.Slice(ciphs, func(i, j int) bool {
		return ciphs[i].freq > ciphs[j].freq
	})

	fmt.Println("Top four results with highest letter frequency Scores:")
	for i := range ciphs[:4] {
		fmt.Printf("XOr'ing Key: %q, Frequency Score: %f, Resulting output: %q\n", ciphs[i].key, ciphs[i].freq, ciphs[i].message)
	}
}

func hexToBytes(s string) (b []byte) {

	str := []byte(s)
	b = make([]byte, hex.DecodedLen(len(str)))
	hex.Decode(b, str)

	return
}

func calcFq(b []byte) (f float64) {

	for _, j := range bytes.ToUpper(b) {
		if v, ok := Fq[j]; ok {
			f += v
		}
	}

	return
}

func xor(b []byte, l byte) (r []byte) {

	for i := range b {
		r = append(r, b[i]^l)
	}

	return
}
