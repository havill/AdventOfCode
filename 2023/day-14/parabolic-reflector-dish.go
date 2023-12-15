package main

import (
	"bufio"
	"fmt"
	"os"
)

func rotate(lines [][]byte) [][]byte {
	rows := len(lines)
	cols := len(lines[0])

	rotated := make([][]byte, cols)
	for i := 0; i < cols; i++ {
		rotated[i] = make([]byte, rows)
		for j := 0; j < rows; j++ {
			rotated[i][j] = lines[rows-1-j][i]
		}
	}
	return rotated
}

func moveNorth(lines [][]byte) {
	for i := 1; i < len(lines); i++ {
		for j := 0; j < len(lines[i]); j++ {
			for k := i - 1; k >= 0; k-- {
				if lines[k][j] == '.' && lines[k+1][j] == 'O' {
					lines[k][j] = 'O'
					lines[k+1][j] = '.'
				}
			}
		}
	}
}

func move(dir byte, lines [][]byte) {
	switch dir {
	case 'N':
		moveNorth(lines)
	case 'S':
		temp := rotate(rotate(lines))
		moveNorth(temp)
		temp = rotate(rotate(temp))
		copy(lines, temp)
	case 'W':
		temp := rotate(lines)
		moveNorth(temp)
		temp = rotate(rotate(rotate(temp)))
		copy(lines, temp)
	case 'E':
		temp := rotate(rotate(rotate(lines)))
		moveNorth(temp)
		temp = rotate(temp)
		copy(lines, temp)
	}
}

func rockTotal(lines [][]byte) int {
	n := 0
	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(lines[i]); j++ {
			if lines[i][j] == 'O' {
				n += len(lines) - i
			}
		}
	}
	return n
}

func cp(lines [][]byte) [][]byte {
	l1 := make([][]byte, len(lines))
	for i, row := range lines {
		l1[i] = append([]byte(nil), row...)
	}
	return l1
}

func goNorth(lines [][]byte) {
	l1 := cp(lines)
	move('N', l1)
	fmt.Println(rockTotal(l1))
}

func key(lines [][]byte) string {
	key := ""
	for i := 0; i < len(lines); i++ {
		key += string(lines[i])
	}
	return key
}

func spinCycle(lines [][]byte) {
	cache := map[string]int{}
	revCache := map[int][][]byte{}
	n := 0
	start := 0
	period := 0
	for ; ; n++ {
		if cache[key(lines)] > 0 {
			if start == 0 {
				start = n
				cache = map[string]int{}
			} else {
				// second time
				period = n - start - 1
				break
			}
		} else {
			cache[key(lines)] = n
			revCache[n] = cp(lines)
		}

		move('N', lines)
		move('W', lines)
		move('S', lines)
		move('E', lines)
	}

	temp := revCache[start+(1000000000-start)%period] // fuck slices
	fmt.Println(rockTotal(temp))
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := [][]byte{} // hope we don't run out of memory
	for scanner.Scan() {
		lines = append(lines, []byte(scanner.Text()))
	}
	goNorth(lines)
	spinCycle(lines)
}
