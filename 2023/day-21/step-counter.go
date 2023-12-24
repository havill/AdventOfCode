package main

import (
	"bufio"
	"fmt"
	"os"
)

func printMap(rocks, plots [][]bool, startX, startY int, occupied rune) {
	for y, row := range rocks {
		for x := range row {
			if x == startX && y == startY {
				fmt.Print("S")
			} else if rocks[y][x] {
				fmt.Print("#")
			} else if plots[y][x] {
				fmt.Print("O")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func countTrue(grid [][]bool) int {
	count := 0
	for _, row := range grid {
		for _, cell := range row {
			if cell {
				count++
			}
		}
	}
	return count
}

func stepCounter(rocks, plots, reached [][]bool, x, y, remaining int) bool {
	walked := false

	if x < 0 || y < 0 || x >= len(rocks[0]) || y >= len(rocks) {
		return false
	}
	if remaining > 0 {
		plots[y][x] = true
		remaining -= 1
		if !rocks[y-1][x] && !plots[y-1][x] { // north
			walked = stepCounter(rocks, plots, reached, x, y-1, remaining) || walked
		}
		if !rocks[y+1][x] && !plots[y+1][x] { // south
			walked = stepCounter(rocks, plots, reached, x, y+1, remaining) || walked
		}
		if !rocks[y][x+1] && !plots[y][x+1] { // east
			walked = stepCounter(rocks, plots, reached, x+1, y, remaining) || walked
		}
		if !rocks[y][x-1] && !plots[y][x-1] { // west
			walked = stepCounter(rocks, plots, reached, x-1, y, remaining) || walked
		}
	} else {
		reached[y][x] = true
	}
	return walked
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var rocks, plots, reached [][]bool
	var startX, startY int

	for j := 0; scanner.Scan(); {
		line := scanner.Text()
		row := make([]bool, len(line))
		for i, char := range line {
			if char == 'S' {
				startX = i
				startY = j
			} else if char == '#' {
				row[i] = true
			} else if char == '.' {
				row[i] = false
			} else {
				fmt.Fprintf(os.Stderr, "Invalid character: %c\n", char)
				os.Exit(1)
			}
		}
		j += 1
		rocks = append(rocks, row)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}

	plots = make([][]bool, len(rocks))
	for i := range rocks {
		plots[i] = make([]bool, len(rocks[i]))
	}
	reached = make([][]bool, len(rocks))
	for i := range rocks {
		reached[i] = make([]bool, len(rocks[i]))
	}

	printMap(rocks, plots, startX, startY, '#')
	stepCounter(rocks, plots, reached, startX, startY, 6)
	fmt.Println(countTrue(reached))
	printMap(rocks, reached, -1, -1, 'O')

}
