package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func isSafe(report []int) bool {
	if len(report) < 2 {
		return true
	}

	increasing := report[1] > report[0]
	for i := 1; i < len(report); i++ {
		diff := int(math.Abs(float64(report[i] - report[i-1])))
		if diff < 1 || diff > 3 {
			return false
		}
		if increasing && report[i] < report[i-1] {
			return false
		}
		if !increasing && report[i] > report[i-1] {
			return false
		}
	}
	return true
}

func main() {
	var reports [][]int

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		numbers := strings.Fields(line)
		var levels []int
		for _, number := range numbers {
			level, err := strconv.Atoi(number)
			if err != nil {
				fmt.Println("Invalid input, please enter integers")
				continue
			}
			levels = append(levels, level)
		}
		reports = append(reports, levels)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}

	safeCount := 0
	for _, report := range reports {
		if isSafe(report) {
			safeCount++
		}
	}
	fmt.Printf("Number of Safe Reports: %d\n", safeCount)
}
