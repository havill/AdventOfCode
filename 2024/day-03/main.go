package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func enableMul() {
	// Placeholder function for enabling multiplication
	fmt.Println("enableMul called")
}

func disableMul() {
	// Placeholder function for disabling multiplication
	fmt.Println("disableMul called")
}

func multiplyFactors(factors string) int {
	parts := strings.Split(factors, ",")
	if len(parts) != 2 {
		fmt.Println("Invalid factors format")
		return 0
	}

	x, err1 := strconv.Atoi(parts[0])
	y, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		fmt.Println("Invalid integers in factors")
		return 0
	}

	return x * y
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanRunes)

	var input string
	for scanner.Scan() {
		input += scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	reMul := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	reDo := regexp.MustCompile(`do\(\)`)
	reDont := regexp.MustCompile(`don't\(\)`)

	state := "ENABLED"
	sum := 0

	for len(input) > 0 {
		dontIndex := reDont.FindStringIndex(input)
		doIndex := reDo.FindStringIndex(input)
		mulIndex := reMul.FindStringSubmatchIndex(input)

		if dontIndex != nil && (doIndex == nil || dontIndex[0] < doIndex[0]) && (mulIndex == nil || dontIndex[0] < mulIndex[0]) {
			disableMul()
			state = "DISABLED"
			input = input[dontIndex[1]:]
		} else if doIndex != nil && (dontIndex == nil || doIndex[0] < dontIndex[0]) && (mulIndex == nil || doIndex[0] < mulIndex[0]) {
			enableMul()
			state = "ENABLED"
			input = input[doIndex[1]:]
		} else if mulIndex != nil && (dontIndex == nil || mulIndex[0] < dontIndex[0]) && (doIndex == nil || mulIndex[0] < doIndex[0]) {
			if state == "ENABLED" {
				factors := input[mulIndex[2]:mulIndex[5]]
				result := multiplyFactors(factors)
				sum += result
				fmt.Printf("Result of mul(%s): %d\n", factors, result)
			}
			input = input[mulIndex[1]:]
		} else {
			break
		}
	}

	fmt.Printf("Sum: %d\n", sum)
}
