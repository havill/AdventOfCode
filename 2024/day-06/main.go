package main

import (
	"bufio"
	"fmt"
	"os"
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
		return lab[y][x] == '#'
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
}
