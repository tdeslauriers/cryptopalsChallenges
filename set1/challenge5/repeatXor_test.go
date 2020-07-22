package main

import "testing"

var key = "ICE"
var txt0 = "Burning 'em, if you ain't quick and nimble I go crazy when I hear a cymbal"

func Test_outputs(t *testing.T) {

	tk := []byte(key)
	tt0 := []byte(txt0)

	t.Logf("\nKey in bytes: %v \nTest text in bytes: %v\n", tk, tt0)
	t.Logf("\nKey in bytes-hex: %x \nTest text in bytes-hex: %x\n", tk, tt0)

}

func Test_xor(t *testing.T) {

	tk := []byte(key)
	tt0 := []byte(txt0)

	x := repeatXor(tt0, tk)
	r := repeatXor(x, tk)

	t.Logf("\nStarting text : %s", tt0)
	t.Logf("\nXor'd text in hex: %x", x)
	t.Logf("\nReturned to clear: %s", r)
}
