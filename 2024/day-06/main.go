package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

type Direction int

const (
	North Direction = 1 << iota // 1 for north
	South                       // 2 for south
	East                        // 4 for east
	West                        // 8 for west
)

func loadMap() [][]rune {
	scanner := bufio.NewScanner(os.Stdin)
	var lab [][]rune

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		row := []rune(line)
		lab = append(lab, row)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return lab
}

func guardPosition(lab [][]rune) (int, int) {
	for y, row := range lab {
		for x, char := range row {
			if char == '^' || char == 'v' || char == '>' || char == '<' {
				return x, y
			}
		}
	}
	return -1, -1 // Return -1, -1 if no guard is found
}

func guardDirection(guard rune) (int, int) {
	switch guard {
	case '^':
		return 0, -1
	case 'v':
		return 0, 1
	case '<':
		return -1, 0
	case '>':
		return 1, 0
	default:
		return 0, 0
	}
}

func isObstacle(lab [][]rune, x, y int) bool {
	if y >= 0 && y < len(lab) && x >= 0 && x < len(lab[y]) {
		return lab[y][x] == '#' || lab[y][x] == 'O'
	}
	return false
}

func isBlocked(lab [][]rune, x, y, dx, dy int) bool {
	newX, newY := x+dx, y+dy
	return isObstacle(lab, newX, newY)
}

func stillOnMap(x, y int, lab [][]rune) bool {
	if y >= 0 && y < len(lab) && x >= 0 && x < len(lab[y]) {
		return true
	}
	return false
}

func traveledRoute(lab [][]rune) int {
	count := 0
	for _, row := range lab {
		for _, char := range row {
			if char == 'X' {
				count++
			}
		}
	}
	return count
}

func moveGuard(x, y, dx, dy int, lab [][]rune) (int, int) {
	lab[y][x] = 'X'
	x += dx
	y += dy
	return x, y
}

func clockwiseTurn(dx, dy int) (int, int) {
	switch {
	case dx == 0 && dy == -1: // North
		return 1, 0 // East
	case dx == 1 && dy == 0: // East
		return 0, 1 // South
	case dx == 0 && dy == 1: // South
		return -1, 0 // West
	case dx == -1 && dy == 0: // West
		return 0, -1 // North
	default:
		return dx, dy // No change if direction is invalid
	}
}

func directionToHex(dir Direction) string {
	return fmt.Sprintf("%X", int(dir))
}

func directionFromDelta(dx, dy int) Direction {
	switch {
	case dx == 0 && dy <= -1:
		return North
	case dx == 0 && dy >= 1:
		return South
	case dx >= 1 && dy == 0:
		return East
	case dx <= -1 && dy == 0:
		return West
	default:
		return 0 // Invalid direction
	}
}

func hexToDirection(hex rune) Direction {
	hex = unicode.ToUpper(hex)
	switch hex {
	case '0':
		return 0
	case '1':
		return North
	case '2':
		return South
	case '3':
		return North | South
	case '4':
		return East
	case '5':
		return North | East
	case '6':
		return South | East
	case '7':
		return North | South | East
	case '8':
		return West
	case '9':
		return North | West
	case 'A':
		return South | West
	case 'B':
		return North | South | West
	case 'C':
		return East | West
	case 'D':
		return North | East | West
	case 'E':
		return South | East | West
	case 'F':
		return North | South | East | West
	default:
		return -1 // Invalid direction
	}
}

func printLab(lab [][]rune) {
	for _, row := range lab {
		fmt.Println(string(row))
	}
	fmt.Println()
}

func main() {
	lab := loadMap()
	printLab(lab)

	x, y := guardPosition(lab)
	dx, dy := guardDirection(lab[y][x])

	originX, originY := x, y
	originDx, originDy := dx, dy
	originLab := make([][]rune, len(lab))
	for i := range lab {
		originLab[i] = make([]rune, len(lab[i]))
		copy(originLab[i], lab[i])
	}

	lab[y][x] = '0'

	for stillOnMap(x, y, lab) {
		for isBlocked(lab, x, y, dx, dy) {
			dx, dy = clockwiseTurn(dx, dy)
		}
		lab[y][x] = 'X'
		x, y = moveGuard(x, y, dx, dy, lab)
	}

	printLab(lab)

	xCount := traveledRoute(lab)
	fmt.Printf("Number of 'X' in the lab: %d\n", xCount)

	// Iterate through the entire lab matrix, resetting x, y, dx, dy, and lab
	for i := range lab {
		for j := range lab[i] {
			x, y = originX, originY
			dx, dy = originDx, originDy
			for k := range lab {
				copy(lab[k], originLab[k])
			}
			if lab[i][j] == '.' {
				lab[i][j] = 'O'
			}
		}
	}

}
