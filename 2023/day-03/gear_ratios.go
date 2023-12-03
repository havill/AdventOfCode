package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type NumberPosition struct {
	StartXPosition int
	EndXPosition   int
	StartYPosition int
	Number         int
}

var positions []NumberPosition

var maxXPosition int = 1000
var maxYPosition int = 1000

func findIntegers(line string, y int) []NumberPosition {
	var positions []NumberPosition
	var startPos, endPos, number int
	reader := strings.NewReader(line)
	for {
		n, err := fmt.Fscanf(reader, "%d", &number)
		if err == io.EOF {
			break
		}
		if n <= 0 {
			reader.ReadByte() // skip non-digit character
			continue
		}
		endPos = int(reader.Size()) - int(reader.Len())        // endPos is the char AFTER the last digit, so no +1 here
		startPos = endPos - len(fmt.Sprintf("%d", number)) + 1 // +1 because we want a 1-based index
		positions = append(positions, NumberPosition{StartXPosition: startPos, EndXPosition: endPos, StartYPosition: y, Number: number})
	}
	return positions
}

func findSymbolsInLine(s string, arr []bool) {
	for i, ch := range s {
		if ch != '.' && (ch < '0' || ch > '9') {
			arr[i+1] = true
		}
	}
}

func CellHasNeighborSymbols(matrix [][]bool, x int, y int) bool {
	for b := -1; b <= 1; b++ {
		for a := -1; a <= 1; a++ {
			if matrix[y+b][x+a] {
				return true
			}
		}
	}
	return false
}

func NumberHasNeighborSymbols(matrix [][]bool, n NumberPosition) int {
	for i := n.StartXPosition; i <= n.EndXPosition; i++ {
		if CellHasNeighborSymbols(matrix, i, n.StartYPosition) {
			return n.Number
		}
	}
	return 0
}

func addAllCellsWithSymbolNeighbors(matrix [][]bool, positions []NumberPosition) int {
	sum := 0

	for _, pos := range positions {
		value := NumberHasNeighborSymbols(matrix, pos)
		if value > 0 {
			// fmt.Printf("%4d at Ln %3d, Col %3d is a part number\n", pos.Number, pos.StartYPosition, pos.StartXPosition)
			sum += value
		}
	}
	return sum
}

func removeSigns(s string) string {
	s = strings.ReplaceAll(s, "+", "#")
	s = strings.ReplaceAll(s, "-", "#")
	return s
}

func debugNumbersAndPositions(positions []NumberPosition) {
	for _, pos := range positions {
		fmt.Printf("Start: %d, End: %d, Y: %d, Number: %d\n", pos.StartXPosition, pos.EndXPosition, pos.StartYPosition, pos.Number)
	}
}

func debugSymbols(matrix [][]bool, maxx int, maxy int) {
	y := 0
	for _, row := range matrix {
		fmt.Printf("%03d: ", y)
		x := 0
		for _, symbol := range row {
			if symbol {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
			x++
			if x > maxx {
				break
			}
		}
		fmt.Println()
		y++
		if y > maxy {
			break
		}
	}
}

func main() {
	var symbols [][]bool

	y := 0

	// Create a 2D array of booleans
	symbols = make([][]bool, maxYPosition)
	for i := range symbols {
		symbols[i] = make([]bool, maxXPosition)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := removeSigns(scanner.Text())
		y++ // y is 1-based
		positions = append(positions, findIntegers(line, y)...)
		findSymbolsInLine(line, symbols[y])
	}
	// debugNumbersAndPositions(positions)
	// debugSymbols(symbols, 140, 140)
	total := addAllCellsWithSymbolNeighbors(symbols, positions)
	fmt.Println(total)
}
