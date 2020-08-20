package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
)

type keylength struct {
	keySize  int
	editDist float64
}

type keylengths []keylength

type cipher struct {
	key   byte
	freq  float64
	block int
}

type ciphers []cipher

type guess struct {
	keysize int
	guess   []byte
	decoded []byte
}

type guesses []guess

func main() {

	// ingest challenge value
	f, err := ioutil.ReadFile("ch6.txt")
	if err != nil {
		panic(err)
	}

	ct, err := base64.StdEncoding.DecodeString(string(f))
	if err != nil {
		panic(err)
	}

	// computes and sorts all key length edit distances, loads keylenth structs
	kls := loadKeyLengths(ct)

	// loads top three edit distance key lengths to slice
	realKeyLengths := make(keylengths, 0, 3)
	for _, j := range kls[:3] {

		realKeyLengths = append(realKeyLengths, j)
	}

	for _, j := range realKeyLengths {

		bg := loadTopGuesses(scoreBlocksSingleCharXor(transpose(ct, j.keySize)), ct, j.keySize)

		fmt.Printf("Key size: %d, \nGuess Key: %q, \nDecoded: %q\n\n", bg[0].keysize, bg[0].guess, bg[0].decoded)
	}

}

// needed this to figure it out: https://golang.org/ref/spec#Arithmetic_operators
func hamm(a, b []byte) (int, error) {

	if len(a) != len(b) {
		return -1, errors.New("compared byte arrays must be the same length")
	}

	hd := 0
	for i := 0; i < len(a); i++ {
		for j := 0; j < 8; j++ {
			mask := byte(1 << uint(j))
			if (a[i] & mask) != (b[i] & mask) {
				hd++
			}
		}
	}

	return hd, nil
}

// loops thru keylengths computing average hamm dist
// loads keylength object
// loads keylengths to array/slice struct
func loadKeyLengths(ct []byte) keylengths {

	var kls keylengths
	for i := 2; i < 42; i++ {

		dnom := len(ct) / (i * 2)
		var num float64
		var ind int

		for j := 0; j < dnom; j++ {

			hm, _ := hamm(ct[ind:ind+i], ct[ind+i:ind+(i*2)])
			num += float64(hm) / float64(i)
			ind = ind + (i * 2)

		}

		nhm := float64(num) / float64(dnom)
		kl := keylength{
			keySize:  i,
			editDist: nhm,
		}

		kls = append(kls, kl)
	}

	sort.SliceStable(kls, func(i, j int) bool {
		return kls[i].editDist < kls[j].editDist
	})

	return kls
}

// breaks cipher text blocks key-sized block
// loads Slice with first index value from each block
// loads next slice with second index value from each block, etc.
func transpose(b []byte, k int) (t [][]byte) {

	n := len(b) / k
	var counter int
	t = make([][]byte, k)
	for i := 0; i < n; i++ {
		for j := range b[counter : counter+k] {

			t[j] = append(t[j], b[counter : counter+k][j])
		}
		counter += k
	}

	return
}

// scores each block slice by xoring with single letter and sorts
func scoreBlocksSingleCharXor(tp [][]byte) []ciphers {

	var allScores []ciphers
	for i := 0; i < len(tp); i++ {
		var ciphs ciphers
		// bytes of 127 most common
		// bytes mapped to byte to frequency in calcFq()
		for j := 0; j < 127; j++ {

			dt := xor(tp[i], byte(j))
			fq := calcFq(dt)
			ci := cipher{
				key:   byte(j),
				freq:  fq,
				block: i,
			}
			ciphs = append(ciphs, ci)

		}

		sort.Slice(ciphs, func(i, j int) bool {
			return ciphs[i].freq > ciphs[j].freq

		})

		allScores = append(allScores, ciphs)
	}

	return allScores
}

func loadTopGuesses(as []ciphers, ct []byte, ks int) (gs guesses) {
	for i := 0; i < 6; i++ {

		//
		var gk []byte
		for _, j := range as {
			gk = append(gk, j[i].key)
		}
		dc := xorWithKey(ct, gk)
		g := guess{
			keysize: ks,
			guess:   gk,
			decoded: dc,
		}
		gs = append(gs, g)
		//
	}
	return
}

// re-used from  challenge 3
func calcFq(b []byte) (f float64) {

	// http://www.fitaly.com/board/domper3/posts/136.html
	var Fq = map[byte]float64{9: 0.0057, 23: 0.0000, 32: 17.1662, 33: 0.0072, 34: 0.2442, 35: 0.0179, 36: 0.0561, 37: 0.0160, 38: 0.0226, 39: 0.2447, 40: 0.2178, 41: 0.2233, 42: 0.0628, 43: 0.0215, 44: 0.7384, 45: 1.3734, 46: 1.5124, 47: 0.1549, 48: 0.5516, 49: 0.4594, 50: 0.3322, 51: 0.1847, 52: 0.1348, 53: 0.1663, 54: 0.1153, 55: 0.1030, 56: 0.1054, 57: 0.1024, 58: 0.4354, 59: 0.1214, 60: 0.1225, 61: 0.0227, 62: 0.1242, 63: 0.1474, 64: 0.0073, 65: 0.3132, 66: 0.2163, 67: 0.3906, 68: 0.3151, 69: 0.2673, 70: 0.1416, 71: 0.1876, 72: 0.2321, 73: 0.3211, 74: 0.1726, 75: 0.0687, 76: 0.1884, 77: 0.3529, 78: 0.2085, 79: 0.1842, 80: 0.2614, 81: 0.0316, 82: 0.2519, 83: 0.4003, 84: 0.3322, 85: 0.0814, 86: 0.0892, 87: 0.2527, 88: 0.0343, 89: 0.0304, 90: 0.0076, 91: 0.0086, 92: 0.0016, 93: 0.0088, 94: 0.0003, 95: 0.1159, 96: 0.0009, 97: 5.1880, 98: 1.0195, 99: 2.1129, 100: 2.5071, 101: 8.5771, 102: 1.3725, 103: 1.5597, 104: 2.7444, 105: 4.9019, 106: 0.0867, 107: 0.6753, 108: 3.1750, 109: 1.6437, 110: 4.9701, 111: 5.7701, 112: 1.5482, 113: 0.0747, 114: 4.2586, 115: 4.3686, 116: 6.3700, 117: 2.0999, 118: 0.8462, 119: 1.3034, 120: 0.1950, 121: 1.1330, 122: 0.0596, 123: 0.0026, 124: 0.0007, 125: 0.0026, 126: 0.0003, 131: 0.0000, 149: 0.6410, 183: 0.0010, 223: 0.0000, 226: 0.0000, 229: 0.0000, 230: 0.0000, 237: 0.0000}

	for _, j := range b {
		if v, ok := Fq[j]; ok {
			f += v
		}
	}

	return
}

// re-used from challenge 3
func xor(b []byte, l byte) (r []byte) {

	for i := range b {
		r = append(r, b[i]^l)
	}

	return
}

func xorWithKey(ct, k []byte) (pt []byte) {

	kc := 0
	for i := range ct {
		pt = append(pt, ct[i]^k[kc])
		if kc >= len(k)-1 {
			kc = 0
		} else {
			kc++
		}
	}

	return
}
