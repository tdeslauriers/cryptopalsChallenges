package main

import (
	"testing"
)

//testing truncated since code is mostly the same from Challenge3.
var testString = "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
var testKey byte = 'X'

func Test_Print(t *testing.T) {

	file := scanFile("./challengeInput.txt")
	t.Logf("%s\n", file[0])
	t.Logf("%x\n", hexToBytes(file[0]))

	// test of known good from previous exercise.
	t.Logf("%q\n", xor(hexToBytes(testString), testKey))
}
