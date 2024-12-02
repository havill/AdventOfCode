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

func problemDampener(report []int) bool {
	//fmt.Printf("testing %v\n", report)
	if isSafe(report) {
		//fmt.Println("...Safe without removing any level")
		return true
	}

	for i := 0; i < len(report); i++ {
		var badLevelRemoved []int
		badLevelRemoved = append(badLevelRemoved, report[:i]...)
		badLevelRemoved = append(badLevelRemoved, report[i+1:]...)
		//fmt.Printf("...testing with %v\n", badLevelRemoved)
		if isSafe(badLevelRemoved) {
			/*
				ending := "th"
				if i == 0 {
					ending = "st"
				} else if i == 1 {
					ending = "nd"
				} else if i == 2 {
					ending = "rd"
				}
				fmt.Printf("... Safe by removing the %d%s level, %d\n", i+1, ending, report[i])
			*/
			return true
		}
	}
	//fmt.Println("...Unsafe regardless of which level is removed")
	return false
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

	safeCount = 0
	for _, report := range reports {
		if problemDampener(report) {
			safeCount++
		}
	}
	fmt.Printf("Number of Safe Reports with Dampening: %d\n", safeCount)
}
