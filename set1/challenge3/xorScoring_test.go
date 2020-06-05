package main

import (
	"sort"
	"testing"
)

// 'Tom' in hex string: short string for tests
var tom = "546f6d"
var let byte = 'D'

func Test_print(t *testing.T) {

	// reviewing generalized outputs
	t.Logf("test hex string back to string: %q", hexToBytes(tom))
	t.Logf("hex string as bytes: %v", []byte(s1))
	t.Logf("Hex to byte, printed as string: %q\n", hexToBytes(s1))

	// keys as bytes
	t.Logf("Keys as bytes: %v", []byte(keys))

	// prints later section of array
	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, j := range s[len(s)-3:] {
		t.Logf("Value: %d", j)
	}

	// map values
	t.Logf("Frequency for 'E': %f", Fq['E'])
	t.Logf("Frequency for 'T': %f", Fq['T'])
}

func Test_xor(t *testing.T) {

	// hex of "D" = "44"
	x := let
	f := xor(hexToBytes(tom), x)
	t.Logf("XOr of 'Tom' with letter D: %v", f)

	// flip it back - re-XOr
	g := xor(f, x)
	t.Logf("Flipping it back: %q", g)
}

func Test_calcFq(t *testing.T) {

	// manual calc of Tom score from Fq
	ltrs := []float64{9.10, 7.68, 2.61}
	tomscore := 9.10 + 7.68 + 2.61
	t.Logf("Tom's manual score of each letter: %f\n", ltrs)
	t.Logf("Tom's manual frequency score: %f\n", tomscore)

	// tom = "546f6d"
	if calcFq(hexToBytes("54")) != ltrs[0] {
		t.Errorf("Expected %f, got %f", ltrs[0], calcFq(hexToBytes("54")))
	}
	if calcFq(hexToBytes("6f")) != ltrs[1] {
		t.Errorf("Expected %f, got %f", ltrs[1], calcFq(hexToBytes("6f")))
	}
	if calcFq(hexToBytes("6d")) != ltrs[2] {
		t.Errorf("Expected %f, got %f", ltrs[2], calcFq(hexToBytes("6d")))
	}
	if calcFq(hexToBytes("546f6d")) != tomscore {
		t.Errorf("Expected %f, got %f", tomscore, calcFq(hexToBytes("546f6d")))
	}

	// score from re-xor of 'Tom'
	x := []byte{16, 43, 41}
	if calcFq(xor(x, let)) != tomscore {
		t.Errorf("Expected xor score %f, got %f", tomscore, calcFq(xor(x, let)))
	}
}

// playground for building main
func Test_loadStruct(t *testing.T) {

	// load one
	m := xor(hexToBytes(s1), let)
	f := calcFq(m)
	out := cipher{
		key:     let,
		message: m,
		freq:    f,
	}
	t.Logf("cipher struct loaded: %x, %s, %v", out.key, out.message, out.freq)
	t.Logf("%s", xor(m, let))

	tKeys := make([]byte, len([]byte(keys)))
	tKeys = []byte(keys)
	t.Logf("keys as bytes: %v", tKeys)

	// load several by loop
	// add to ciphers array
	var ci ciphers
	for _, j := range tKeys[:] {
		t.Logf("Test Keys: %v", j)
		b := xor(hexToBytes(s1), j)
		c := calcFq(b)
		o := cipher{
			key:     j,
			message: b,
			freq:    c,
		}
		t.Logf("Test loop, load cipher: %q, %s, %v", o.key, o.message, o.freq)
		ci = append(ci, o)
	}
	t.Logf("ciphers array: %v", ci)
	for i := range ci {
		t.Logf("Ciphers index %d: %q, %f, %q", i, ci[i].key, ci[i].freq, ci[i].message)
	}

	// test sort
	sort.Slice(ci, func(i, j int) bool {
		return ci[i].freq > ci[j].freq
	})
	for i := range ci[:4] {
		t.Logf("Ciphers(sorted) index %d: %q, %f, %q", i, ci[i].key, ci[i].freq, ci[i].message)
	}
}

// can be done without go funcs, this improves performance
func Test_WaitGroups(t *testing.T) {

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
	for i := range ciphs[:4] {
		t.Logf("WaitGroup ciphers(sorted) index %d: %q, %f, %q", i, ciphs[i].key, ciphs[i].freq, ciphs[i].message)
	}
}
