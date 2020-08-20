package main

import (
	"encoding/base64"
	"io/ioutil"
	"testing"
)

var t1 = "this is a test"
var t2 = "wokka wokka!!!"
var t3 = "monekywrench"
var hammD = 37

var testText = `Three Rings for the Elven-kings under the sky,
Seven for the Dwarf-lords in their halls of stone,
Nine for Mortal Men doomed to die,
One for the Dark Lord on his dark throne
In the Land of Mordor where the Shadows lie.
One Ring to rule them all, One Ring to find them,
One Ring to bring them all, and in the darkness bind them,
In the Land of Mordor where the Shadows lie.`

var testKey = "Sauron's Fiery Eye"

func Test_hamm(t *testing.T) {

	c := []byte(t1)
	d := []byte(t2)
	e := []byte(t3)

	d1, err := hamm(c, d)
	if err != nil {
		t.Logf("Expected err == nil: got %v", err)
	}
	t.Logf("Hamming distance == %d", d1)

	// getting it wrong on purpose
	_, err = hamm(c, e)
	if err != nil {
		t.Logf("Expected err != nil: got %v", err)
	}

	// slight score difference
	f := []byte("this is a test")
	g := []byte("tokka wokka!!!")

	d1, err = hamm(f, g)
	if err != nil {
		t.Logf("Expected err == nil: got %v", err)
	}
	t.Logf("Hamming distance == %d", d1)
}

func Test_outputs(t *testing.T) {

	// reads to bytes
	f, err := ioutil.ReadFile("ch6.txt")
	if err == nil {

		t.Logf("%s", f)
	}

	// "tomtom"
	if dt, err := base64.StdEncoding.DecodeString("dG9tdG9tCg=="); err == nil {

		t.Logf("%s", dt)
	}

}

func Test_xorWithKey(t *testing.T) {

	tct := []byte(testText)
	tk := []byte(testKey)

	x1 := xorWithKey(tct, tk)
	xb := xorWithKey(x1, tk)

	if string(tct) != string(xb) {

		t.Errorf("Expected return Xor to be 'One ring...': was %q", xb)
	}
	t.Logf("'One ring...' XOr'd: %x", x1)
	t.Logf("'One ring...' XOr'd back: %q", xb)
}

// new testing from overhaul.
func Test_loadKeyLengths(t *testing.T) {

	testScores := xorWithKey([]byte(testText), []byte(testKey))
	possibleKeys := loadKeyLengths(testScores)

	for _, j := range possibleKeys {

		t.Logf("Key size: %d, Edit distance: %f", j.keySize, j.editDist)
	}
}

func Test_transpose(t *testing.T) {

	tr0 := transpose([]byte("GollumGollumGollumGollumGollumGollumGollumGollumGollum"), len(testKey))
	t.Logf("%s", tr0)

	tr1 := transpose([]byte("SamSamSamSamSamSamSamSamSamSam"), 3)
	t.Logf("%s", tr1)

}

func Test_scoreBlockSingleCharXor(t *testing.T) {

	tr0 := transpose((xorWithKey([]byte(testText), []byte(testKey))), len(testKey))
	tScores := scoreBlocksSingleCharXor(tr0)
	for _, j := range tScores {
		for _, k := range j[:6] {

			t.Logf("Block: %d, Key: %q, Frequency: %f\n", k.block, k.key, k.freq)
		}
	}
}

func Test_loadTop6Guesses(t *testing.T) {

	tr0 := transpose((xorWithKey([]byte(testText), []byte(testKey))), len(testKey))
	tScores := scoreBlocksSingleCharXor(tr0)
	bgs := loadTopGuesses(tScores, xorWithKey([]byte(testText), []byte(testKey)), len(testKey))
	for _, j := range bgs {

		t.Logf("Key size: %d, Guess Key: %q, Decoded: %q\n", j.keysize, j.guess, j.decoded[:100])
	}
}

func Test_full(t *testing.T) {

	// // gut check:
	// gc := loadKeyLengths([]byte(t1 + t2))
	// for _, j := range gc {

	// 	t.Logf("Key size: %d, Edit Dist: %f", j.keySize, j.editDist)
	// }

	encoded := xorWithKey([]byte(testText), []byte(testKey))
	keylens := loadKeyLengths(encoded)
	for _, j := range keylens[:3] {
		t.Logf("Full Test - KeySize: %d, Edit Distance: %f\n", j.keySize, j.editDist)
	}

	tp := transpose(encoded, keylens[0].keySize)
	sb := scoreBlocksSingleCharXor(tp)
	bg := loadTopGuesses(sb, encoded, keylens[0].keySize)
	for _, j := range bg {
		t.Logf("Key size: %d, Guess Key: %q, Decoded: %q\n", j.keysize, j.guess, j.decoded[:100])
	}
}
