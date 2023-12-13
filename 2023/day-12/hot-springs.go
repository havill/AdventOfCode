package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func splitString(input string) (string, string) {
	parts := strings.Split(input, " ")
	left := ""
	right := ""

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.ContainsAny(part, ".#?") {
			left += part
		} else {
			right += part
		}
	}

	return left, right
}

func stringToIntSlice(input string) ([]int, error) {
	parts := strings.Split(input, ",")
	result := make([]int, len(parts))

	for i, part := range parts {
		part = strings.TrimSpace(part)
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		result[i] = num
	}

	return result, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please enter some lines of text. Press CTRL+D to end input.")

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("You entered: ", line)
		left, right := splitString(line)
		fmt.Println("Left: ", left)
		fmt.Println("Right: ", right)
		brokenGroups, err := stringToIntSlice(right)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		fmt.Println("Broken groups: ", brokenGroups)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
