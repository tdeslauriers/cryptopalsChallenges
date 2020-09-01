package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	file, err := os.Open("./ch8.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		ct := scanner.Text()
		s := len(ct) / 16
		for i := 0; i < s; i++ {

			for j := 0; j < s; j++ {

				if j != i {

					if ct[i*16:(i+1)*16] == ct[j*16:(j+1)*16] {

						fmt.Printf("These two indexes are the same 16 byte blocks: %d, %d\n", i, j)
						fmt.Printf("This cipher-text is likely an EBC cipher: %v\n", ct)
					}
				}
			}
		}
	}
}
