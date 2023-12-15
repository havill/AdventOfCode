package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func calculateHash(s string) int {
	hash := 0
	for _, c := range s {
		hash = (hash + int(c)) * 17 % 256
	}
	return hash
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	// Remove all whitespace
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "\n", "")

	// Split the input into steps
	steps := strings.Split(input, ",")

	// Calculate the sum of the hash values
	sum := 0
	for _, step := range steps {
		sum += calculateHash(step)
	}

	fmt.Println(sum)
}
