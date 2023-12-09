package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func diffSlice(slice []int) []int {
	var diff []int
	for i := 0; i < len(slice)-1; i++ {
		diff = append(diff, slice[i+1]-slice[i])
	}
	return diff
}

func allZeros(slice []int) bool {
	for _, value := range slice {
		if value != 0 {
			return false
		}
	}
	return true
}

func extrapolateForward(history [][]int) int {
	var sum int = 0

	last := len(history) - 1
	history[last] = append(history[last], 0)

	for i := last; i > 0; i-- {
		previous := i - 1
		addend1 := history[i][len(history[i])-1]
		addend2 := history[previous][len(history[previous])-1]
		sum = addend1 + addend2
		history[previous] = append(history[previous], sum)
	}
	return sum
}

func extrapolateBackwards(history [][]int) int {
	var difference int = 0

	last := len(history) - 1
	history[last] = append([]int{0}, history[last]...)

	for i := last; i > 0; i-- {
		previous := i - 1
		subtrahend := history[i][0]
		minuend := history[previous][0]
		difference = minuend - subtrahend
		history[previous] = append([]int{difference}, history[previous]...)
	}
	return difference
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var histories [][][]int
	var total int

	for scanner.Scan() {
		dataset := scanner.Text()
		numbers := strings.Fields(dataset)
		if len(numbers) < 1 {
			fmt.Println("Invalid input. Please enter at least one integer.")
			continue
		}

		var ints []int
		var extrapolation [][]int
		for _, number := range numbers {
			n, err := strconv.Atoi(number)
			if err != nil {
				fmt.Printf("Invalid number: %s\n", number)
				continue
			}
			ints = append(ints, n)
		}

		extrapolation = append(extrapolation, ints)
		histories = append(histories, extrapolation)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	for i := 0; i < len(histories); i++ {
		for j := 0; j < len(histories[i]); j++ {
			if allZeros(histories[i][j]) {
				break
			}
			histories[i] = append(histories[i], diffSlice(histories[i][j]))
		}
	}

	total = 0
	for _, i := range histories {
		total += extrapolateForward(i)
	}
	fmt.Println(total)

	total = 0
	for _, i := range histories {
		total += extrapolateBackwards(i)
	}
	fmt.Println(total)
}
