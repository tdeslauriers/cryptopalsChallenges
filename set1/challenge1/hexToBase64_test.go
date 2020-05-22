package main

import (
	"testing"
)

var testinput []byte = []byte("49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d")

// must convert string to bytes
func Test_easy(t *testing.T) {

	t.Log(testinput)
}

// must decode hex
func Test_DecodeHex(t *testing.T) {

	if x, err := DecodeHex(testinput); err != nil {
		t.Log(err)
		t.Errorf("Decode failed, returned %v", x)
	} else {
		t.Logf("Decoded successfully: %v", x)
	}
}

// must encode to base64
func Test_Base64Encode(t *testing.T) {

	x, err := DecodeHex(testinput)
	if err != nil {
		t.Logf("Failed to decode hex: %s", err)
	}
	y := Base64Encode(x)
	t.Logf("Base64 encoding outcome: %v", y)
}
