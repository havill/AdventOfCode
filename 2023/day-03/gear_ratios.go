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

func findDigitsInLine(s string, arr []bool) {
	for i, ch := range s {
		if ch >= '0' && ch <= '9' {
			arr[i+1] = true
		}
	}
}

func findGearsInLine(s string, arr []bool) {
	for i, ch := range s {
		if ch == '*' {
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

func getPartsSurroundingGear(positions []NumberPosition, digits [][]bool, x int, y int) (int, int) {
	factor1, factor2 := 0, 0

	for b := y - 1; b <= y+1; b++ {
		for a := x - 1; a <= x+1; a++ {
			if a == x && b == y {
				continue
			}
			for _, pos := range positions {
				if pos.StartXPosition <= a && pos.EndXPosition >= a && pos.StartYPosition == b {
					if factor1 == 0 {
						// fmt.Println("Found factor 1 at Ln ", b, " Col ", a)
						factor1 = pos.Number
					} else if factor2 == 0 {
						// fmt.Println("Found factor 2 at Ln ", b, " Col ", a)
						factor2 = pos.Number
					}
					for digits[b][a] {
						a++ // skip the digit
					}
					// fmt.Println("Skipping to Ln ", b, " Col ", a)
				}
			}
		}
	}
	return factor1, factor2
}

func scanGearsForParts(gears, digits [][]bool, maxx int, maxy int) int {
	sum := 0
	y := 0
	for _, row := range gears {
		x := 0
		for _, gear := range row {
			if gear {
				factor1, factor2 := getPartsSurroundingGear(positions, digits, x, y)
				// fmt.Print("Gear at ", x, ",", y, " has parts ", factor1, " and ", factor2, "\n")
				sum += factor1 * factor2
			}
			x++
			if x > maxx {
				break
			}
		}
		y++
		if y > maxy {
			break
		}
	}
	return sum
}

func debugNumbersAndPositions(positions []NumberPosition) {
	for _, pos := range positions {
		fmt.Printf("Start: %d, End: %d, Y: %d, Number: %d\n", pos.StartXPosition, pos.EndXPosition, pos.StartYPosition, pos.Number)
	}
}

func debugMatrix(matrix [][]bool, maxx int, maxy int) {
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
	var symbols, gears, digits [][]bool

	y := 0

	// Create three 2D arrays of booleans
	symbols = make([][]bool, maxYPosition)
	for i := range symbols {
		symbols[i] = make([]bool, maxXPosition)
	}
	gears = make([][]bool, maxYPosition)
	for i := range gears {
		gears[i] = make([]bool, maxXPosition)
	}
	digits = make([][]bool, maxYPosition)
	for i := range digits {
		digits[i] = make([]bool, maxXPosition)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := removeSigns(scanner.Text())
		y++ // y is 1-based
		positions = append(positions, findIntegers(line, y)...)
		findSymbolsInLine(line, symbols[y])
		findGearsInLine(line, gears[y])
		findDigitsInLine(line, digits[y])
	}
	// debugNumbersAndPositions(positions)
	// debugMatrix(symbols, 140, 140)
	// debugMatrix(gears, 10, 10)
	// debugMatrix(digits, 10, 10)

	total := addAllCellsWithSymbolNeighbors(symbols, positions)
	fmt.Println(total)

	total = scanGearsForParts(gears, digits, maxXPosition, maxYPosition)
	fmt.Println(total)
}
