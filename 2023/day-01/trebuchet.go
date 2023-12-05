package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	sum := 0
	for scanner.Scan() {
		line := scanner.Text()
		firstDigit := ""
		for _, r := range line {
			if unicode.IsDigit(r) {
				firstDigit = string(r)
				break
			}
		}
		lastDigit := ""
		for i := len(line) - 1; i >= 0; i-- {
			if unicode.IsDigit(rune(line[i])) {
				lastDigit = string(line[i])
				break
			}
		}
		calibrationValue, _ := strconv.Atoi(firstDigit + lastDigit)
		sum += calibrationValue
	}
	fmt.Println(sum)
}
