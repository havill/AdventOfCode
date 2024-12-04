package main

import (
	"bufio"
	"fmt"
	"os"
)

func xmasSearcher(wordSearch [][]rune) int {
	directions := [][2]int{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1}, // up, down, left, right
		{-1, -1}, {-1, 1}, {1, -1}, {1, 1}, // up-left, up-right, down-left, down-right
	}

	rows := len(wordSearch)
	cols := len(wordSearch[0])
	count := 0

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if wordSearch[i][j] == 'X' {
				for _, dir := range directions {
					x, y := i+dir[0], j+dir[1]
					if x >= 0 && x < rows && y >= 0 && y < cols && wordSearch[x][y] == 'M' {
						x2, y2 := x+dir[0], y+dir[1]
						if x2 >= 0 && x2 < rows && y2 >= 0 && y2 < cols && wordSearch[x2][y2] == 'A' {
							x3, y3 := x2+dir[0], y2+dir[1]
							if x3 >= 0 && x3 < rows && y3 >= 0 && y3 < cols && wordSearch[x3][y3] == 'S' {
								count++
							}
						}
					}
				}
			}
		}
	}

	return count
}

func crossMasSearcher(wordSearch [][]rune) int {
	rows := len(wordSearch)
	cols := len(wordSearch[0])
	count := 0

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if wordSearch[i][j] == 'A' {
				// Check upper-left and lower-right
				if i > 0 && j > 0 && i < rows-1 && j < cols-1 {
					if wordSearch[i-1][j-1] == 'M' && wordSearch[i+1][j+1] == 'S' || wordSearch[i-1][j-1] == 'S' && wordSearch[i+1][j+1] == 'M' {
						// Check upper-right and lower-left
						if wordSearch[i-1][j+1] == 'M' && wordSearch[i+1][j-1] == 'S' || wordSearch[i-1][j+1] == 'S' && wordSearch[i+1][j-1] == 'M' {
							count++
						}
					}
				}
			}
		}
	}

	return count
}

func main() {
	var wordSearch [][]rune
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter lines of text (Ctrl+D to end):")
	for scanner.Scan() {
		line := scanner.Text()
		var chars []rune
		for _, char := range line {
			chars = append(chars, char)
		}
		wordSearch = append(wordSearch, chars)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	// Print the 2D array
	for _, line := range wordSearch {
		for _, char := range line {
			fmt.Printf("%c ", char)
		}
		fmt.Println()
	}

	count := xmasSearcher(wordSearch)
	fmt.Println("Number of 'XMAS' sequences found:", count)
}
