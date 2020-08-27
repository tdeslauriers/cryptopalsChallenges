package main

import (
	"crypto/aes"
	"encoding/base64"
	"io/ioutil"
	"testing"
)

func Test_inputLoad(t *testing.T) {

	file, err := ioutil.ReadFile("ch7.txt")
	if err != nil {
		t.Log(err)
	}

	ct, err := base64.RawStdEncoding.DecodeString(string(file))
	if err != nil {
		t.Log(err)
	}

	t.Logf("Printing to 50 bytes to hex to ensure file loaded:\n\n %x", ct[:50])
}

func Test_bytesDivisable(t *testing.T) {

	t.Logf("Key byte length: %d\n", len(key))

	file, err := ioutil.ReadFile("ch7.txt")
	if err != nil {
		t.Log(err)
	}

	ct, err := base64.RawStdEncoding.DecodeString(string(file))
	if err != nil {
		t.Log(err)
	}

	t.Logf("Byte count of cypher text: %d\n", len(ct))
	if !bytesDivisable(len(ct), len(key)) {
		t.Logf("Expected bytes divisible: got: %v", bytesDivisable(len(ct), len(key)))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		t.Log(err.Error())
	}
	if !bytesDivisable(len(ct), block.BlockSize()) {
		t.Logf("Expected bytes divisible: got: %v", bytesDivisable(len(ct), len(key)))
	}
}

func Test_decryptAesEcb(t *testing.T) {

	file, err := ioutil.ReadFile("ch7.txt")
	if err != nil {
		t.Log(err)
	}

	ct, err := base64.RawStdEncoding.DecodeString(string(file))
	if err != nil {
		t.Log(err)
	}

	pt := decryptAesEcb(ct, key)
	t.Logf("\nDecrypted text: \n\n%q", pt)
}
