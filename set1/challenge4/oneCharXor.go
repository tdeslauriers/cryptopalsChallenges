package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

var keys = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// http://pi.math.cornell.edu/~mec/2003-2004/cryptography/subs/frequencies.html
var Fq = map[byte]float64{'E': 12.02, 'T': 9.10, 'A': 8.12, 'O': 7.68, 'I': 7.31, 'N': 6.95, 'S': 6.28, 'R': 6.02, 'H': 5.92, 'D': 4.32, 'L': 3.98, 'U': 2.88, 'C': 2.71, 'M': 2.61, 'F': 2.30, 'Y': 2.11, 'W': 2.09, 'G': 2.03, 'P': 1.82, 'B': 1.49, 'V': 1.11, 'K': 0.69, 'X': 0.17, 'Q': 0.11, 'J': 0.10, 'Z': 0.07}

type cipher struct {
	key          byte
	startMessage []byte
	xorMessage   []byte
	freq         float64
}

type ciphers []cipher

var wg sync.WaitGroup

func main() {

	flag.Parse()
	file := scanFile(flag.Arg(0))
	k := []byte(keys)
	var ciphs ciphers

	wg.Add(len(file))
	c := make(chan cipher)

	for _, line := range file {

		go func(str string) {

			for _, j := range k {

				h := hexToBytes(str)
				x := xor(h, j)
				f := calcFq(x)
				ci := cipher{
					key:          j,
					startMessage: h,
					xorMessage:   x,
					freq:         f,
				}
				c <- ci
			}
			wg.Done()
		}(line)

	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for n := range c {
		ciphs = append(ciphs, n)
	}

	sort.Slice(ciphs, func(i, j int) bool {
		return ciphs[i].freq > ciphs[j].freq
	})

	for i := range ciphs[:5] {
		fmt.Printf("XOr'ing Results:\n "+
			"\tKey: %q\n "+
			"\tFrequency Score: %f\n "+
			"\tStarting hex: %x\n "+
			"\tResulting output: %q\n\n",
			ciphs[i].key, ciphs[i].freq, ciphs[i].startMessage, ciphs[i].xorMessage)
	}
}

func scanFile(s string) (lines []string) {

	f, err := os.Open(s)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	fn := bufio.NewScanner(f)

	for fn.Scan() {
		lines = append(lines, fn.Text())
	}

	return
}

func hexToBytes(s string) (b []byte) {

	str := []byte(s)
	b = make([]byte, hex.DecodedLen(len(str)))
	hex.Decode(b, str)

	return
}

func xor(b []byte, l byte) (r []byte) {

	for i := range b {
		r = append(r, b[i]^l)
	}

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
